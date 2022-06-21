# Terra Testnet Faucet

Paloma Testnet Faucet is a client tool that allows anyone to easily request a nominal amount of Terra or Luna assets for testing purposes. This app needs to be deployed on a Paloma testnet full node.

**WARNING**: Tokens recieved over the faucet are not real assets and have no market value.

This faucet implementation is a fork of the [Cosmos Faucet](https://github.com/cosmos/faucet).

## Get tokens on Terra testnets

Using the testnets is really easy. Simply go to https://github.com/palomachain/paloma and follow the instructions to get GRAINs. 

## Usage

Build the docker image.

```bash
docker build -t faucet .
```

Run it with the mnemonic and recaptcha key as env vars.

```bash
docker run -p 3000:3000 \
    -e MNEMONIC=$MY_MNEMONIC \
    -e RECAPTCHA_KEY=$RECAPTCHA_KEY \
    -e PORT=8080 \  # default to 3000
    faucet
```
