module github.com/terra-project/faucet

require (
	github.com/rs/cors v1.8.2
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/tmlibs v0.0.0-20180607034639-640af0205d98
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
)

require github.com/dpapathanasiou/go-recaptcha v0.0.0-20190121160230-be5090b17804

require (
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

go 1.18
