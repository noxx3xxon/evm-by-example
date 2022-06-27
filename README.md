# EVM BY EXAMPLE

## Installation

- Install Golang https://go.dev/doc/install

```
git clone git@github.com:noxx3xxon/evm-by-example.git
cd evm-by-example
go mod tidy
```

The repo queries blockchain nodes to simulate examples, it does this via an environment variable ETH_RPC_URL. Use a provider like [Alchemy](https://alchemy.com/) to get a node RPC endpoint.

```
export ETH_RPC_URL=https://eth-mainnet.alchemyapi.io/v2/[API_KEY]
```

## Structure

This repository extracts sections of the [Geth codebase](https://github.com/ethereum/go-ethereum) so that they can be run in isolation. They are designed to be examples that improve our understanding of how the EVM is implemented in an Ethereum Client. The sub folders represent different areas of the EVM that have been isolated.

To run the code associated with a section

```
go run [SECTION_NAME]/main.go
```

Each individual section has its own README which will provide more details on what that section covers.