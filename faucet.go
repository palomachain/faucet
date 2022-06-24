package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"

	"github.com/rs/cors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tmlibs/bech32"
	"github.com/tomasen/realip"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/core/v2/app"
	"github.com/terra-money/core/v2/app/params"
	//"github.com/tendermint/tendermint/crypto"
)

var recaptchaKey string
var palomad string
var port string
var chainID string
var rpcUrl string

var bankAddress string
var mtx sync.Mutex
var isClassic bool

const ( // new core hasn't these yet.
	MicroUnit              = int64(1e6)
	fullFundraiserPath     = "m/44'/330'/0'/0/0"
	accountAddresPrefix    = "paloma"
	accountPubKeyPrefix    = "palomapub"
	validatorAddressPrefix = "palomavaloper"
	validatorPubKeyPrefix  = "palomavaloperpub"
	consNodeAddressPrefix  = "palomavalcons"
	consNodePubKeyPrefix   = "palomavalconspub"
)

var amountTable = map[string]int64{
	"ugrain": 10 * MicroUnit,
}

const (
	requestLimitSecs = 30
	mnemonicVar      = "MNEMONIC"
	privkeyVar       = "PRIV_KEY"
	recaptchaKeyVar  = "RECAPTCHA_KEY"
	portVar          = "PORT"
	lcdUrlVar        = "LCD_URL"
	chainIDVar       = "CHAIN_ID"
)

// Claim wraps a faucet claim
type Claim struct {
	Address  string `json:"address"`
	Response string `json:"response"`
	Denom    string `json:"denom"`
}

// Coin is the same as sdk.Coin
type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

func newCodec() *params.EncodingConfig {
	ec := app.MakeEncodingConfig()

	config := sdk.GetConfig()
	config.SetCoinType(app.CoinType)
	config.SetFullFundraiserPath(fullFundraiserPath)
	config.SetBech32PrefixForAccount(accountAddresPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()

	return &ec
}

type CoreCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type BalanceResponse struct {
	Balance CoreCoin `json:"balance"`
}

func getBalance(address string) (amount int64) {
	cmd := exec.Command(
		palomad,
		"--node", rpcUrl,
		"q", "bank", "balances",
		"--output", "json",
		"--denom", "ugrain",
		"--chain-id", chainID,
		address,
	)
	res, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("error getting the balance information")
		fmt.Println(string(res))
		panic(err)
	}
	var balance struct {
		Amount string `json:"amount"`
	}

	err = json.Unmarshal(res, &balance)
	if err != nil {
		panic(err)
	}
	amount, err = strconv.ParseInt(balance.Amount, 10, 64)
	if err != nil {
		panic(err)
	}

	return amount
}

func parseRegexp(regexpStr string, target string) (data string) {
	// Capture seqeunce string from json
	r := regexp.MustCompile(regexpStr)
	groups := r.FindStringSubmatch(string(target))

	if len(groups) != 2 {
		os.Exit(1)
	}

	// Convert sequence string to int64
	data = groups[1]
	return
}

// RequestLog stores the Log of a Request
type RequestLog struct {
	Coins     []Coin    `json:"coin"`
	Requested time.Time `json:"updated"`
}

func (requestLog *RequestLog) dripCoin(denom string) error {
	amount := amountTable[denom]

	// try to update coin
	for idx, coin := range requestLog.Coins {
		if coin.Denom == denom {
			if (requestLog.Coins[idx].Amount + amount) > amountTable[denom]*2 {
				return errors.New("amount limit exceeded")
			}

			requestLog.Coins[idx].Amount += amount
			return nil
		}
	}

	// first drip for denom
	requestLog.Coins = append(requestLog.Coins, Coin{Denom: denom, Amount: amount})
	return nil
}

func checkAndUpdateLimit(db *leveldb.DB, account []byte, denom string) error {
	address, _ := bech32.ConvertAndEncode("paloma", account)

	if getBalance(address) >= amountTable[denom]*2 {
		return errors.New("amount limit exceeded")
	}

	var requestLog RequestLog

	logBytes, _ := db.Get(account, nil)
	now := time.Now()

	if logBytes != nil {
		jsonErr := json.Unmarshal(logBytes, &requestLog)
		if jsonErr != nil {
			return jsonErr
		}

		// check interval limt
		intervalSecs := now.Sub(requestLog.Requested).Seconds()
		if intervalSecs < requestLimitSecs {
			return errors.New("please wait a while for another tap")
		}

		// reset log if month was changed
		if requestLog.Requested.Month() != now.Month() {
			requestLog.Coins = []Coin{}
		}

		// check amount limit
		dripErr := requestLog.dripCoin(denom)
		if dripErr != nil {
			return dripErr
		}
	}

	// update requested time
	requestLog.Requested = now
	logBytes, _ = json.Marshal(requestLog)
	updateErr := db.Put(account, logBytes, nil)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func createGetCoinsHandler(db *leveldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, err.(error).Error(), 400)
			}
		}()

		var claim Claim

		// decode JSON response from front end
		decoder := json.NewDecoder(request.Body)
		decoderErr := decoder.Decode(&claim)

		if decoderErr != nil {
			panic(decoderErr)
		}

		amount, ok := amountTable[claim.Denom]

		if !ok {
			panic(fmt.Errorf("invalid denom; %v", claim.Denom))
		}

		// make sure address is bech32
		readableAddress, decodedAddress, decodeErr := bech32.DecodeAndConvert(claim.Address)
		if decodeErr != nil {
			panic(decodeErr)
		}
		// re-encode the address in bech32
		encodedAddress, encodeErr := bech32.ConvertAndEncode(readableAddress, decodedAddress)
		if encodeErr != nil {
			panic(encodeErr)
		}

		// make sure captcha is valid
		clientIP := realip.FromRequest(request)
		captchaResponse := claim.Response
		captchaPassed, captchaErr := recaptcha.Confirm(clientIP, captchaResponse)
		if captchaErr != nil {
			panic(captchaErr)
		}
		if !captchaPassed {
			err := errors.New("captcha failed, please refresh page and try again")
			panic(err)
		}
		// send the coins!

		// Limiting request speed
		limitErr := checkAndUpdateLimit(db, decodedAddress, claim.Denom)
		if limitErr != nil {
			panic(limitErr)
		}

		mtx.Lock()
		defer mtx.Unlock()

		fmt.Println(time.Now().UTC().Format(time.RFC3339), "req", clientIP, encodedAddress, amount, claim.Denom)

		cmd := exec.Command(
			palomad,
			"--node", rpcUrl,
			"tx", "bank", "send",
			"-y",
			"--broadcast-mode", "block",
			"--chain-id", chainID,
			"--fees", "200000ugrain",
			bankAddress,
			encodedAddress,
			fmt.Sprintf("%d%s", amount, claim.Denom),
		)
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("error running a command: %s\n", err)
			fmt.Println("output:")
			fmt.Println(string(output))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"amount": %v}`, amount)
		return

	}
}

func main() {

	bankAddress = os.Getenv("BANK_ADDR")
	if bankAddress == "" {
		panic("BANK_ADDR variable is required")
	}

	palomad = os.Getenv("PALOMA_CMD")
	if palomad == "" {
		panic("PALOMA_CMD variable is required")
	}

	rpcUrl = os.Getenv("NODE_RPC_URL")
	if palomad == "" {
		panic("NODE_RPC_URL variable is required")
	}

	recaptchaKey = os.Getenv(recaptchaKeyVar)

	if recaptchaKey == "" {
		panic("RECAPTCHA_KEY variable is required")
	}

	port = os.Getenv(portVar)

	if port == "" {
		port = "3000"
	}

	chainID = os.Getenv(chainIDVar)

	if chainID == "" {
		panic("CHAIN_ID variable is required")
	}

	db, err := leveldb.OpenFile("db/ipdb", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	recaptcha.Init(recaptchaKey)

	// Application server.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	mux.HandleFunc("/claim", createGetCoinsHandler(db))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://faucet.palomachain.com", "http://localhost", "localhost", "http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatal("failed to start server", err)
	}
}
