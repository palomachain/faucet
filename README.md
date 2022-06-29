# Paloma Testnet Faucet

Paloma Testnet Faucet is a client tool that allows anyone to easily request a nominal amount of Paloma GRAINs and Paloma assets for testing purposes. This app needs to be deployed on a Paloma testnet full node.

**WARNING**: Tokens recieved over the faucet are not real assets and have no market value.

This faucet implementation is a fork of the [Cosmos Faucet](https://github.com/cosmos/faucet).

## Get tokens on Paloma testnets

Using the testnets is really easy. Simply go to https://github.com/palomachain/paloma and follow the instructions to get GRAINs. 

## Usage
For FE development see [Readme in frontend folder](https://github.com/palomachain/faucet/tree/main/frontend)


Build the docker image.

```bash
docker build -t faucet .
```

Run it with the mnemonic and recaptcha key as env vars.

```bash
docker run -p 3000:3000 \
    -e MNEMONIC="$(cat mnemonic.txt)" \
    -e RECAPTCHA_KEY=potato \
    -e PORT=8080 \
    -e LCD_URL=http://165.232.91.129:1317 \
    -e CHAIN_ID=paloma \
    faucet
```
