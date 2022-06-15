# Terra Testnet Faucet

![banner](./terra-faucet.png)

Terra Testnet Faucet is a client tool that allows anyone to easily request a nominal amount of Terra or Luna assets for testing purposes. This app needs to be deployed on a Terra testnet full node, because it relies on using the `terracli` command to send tokens.

**WARNING**: Tokens recieved over the faucet are not real assets and have no market value.

This faucet implementation is a fork of the [Cosmos Faucet](https://github.com/cosmos/faucet).

## Get tokens on Terra testnets

Using the testnets is really easy. Simply go to https://faucet.terra.money, chose your network and input your testnet address. 

## Usage

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
