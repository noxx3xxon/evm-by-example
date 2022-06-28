# BLOOM FILTERS

This section covers bloom filters, a key data structure used to create the logsBloom in both the block header and transaction receipts.

See the "Bloom Filter" section in the [EVM Deep Dives Part 6](https://noxx.substack.com/p/evm-deep-dives-the-path-to-shadowy-16e) article for more details.

To run the code associated with this section first ensure the ETH_RPC_URL env var is set then run main.go

```
export ETH_RPC_URL=https://eth-mainnet.alchemyapi.io/v2/[API_KEY]
go run bloom/main.go
```

## Example - Block 15001871

This code section run through the bloom filters contained in [Block 15001871](https://etherscan.io/block/15001871). The following print statements have been added to take you from block header to the specific bits that are flipped within the bloom filter.

- Prints Block Header for [block 15001871](https://etherscan.io/block/15001871)
- Prints Transaction Data for first Tx in Block 15001871 [0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58](https://etherscan.io/tx/0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58)
- Prints Transaction Receipt for [0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58](https://etherscan.io/tx/0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58)
- Prints Event Logs for Transaction Receipt [0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58](https://etherscan.io/tx/0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58)
- Prints out bit flip indexes in bloom filter caused by the Event Log of Transaction Receipt [0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58](https://etherscan.io/tx/0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58)
- Prints out Bloom Filter in hexadecimal
- Prints out Bloom filter in binary
- Prints out specific bytes within Bloom Filter that were referenced in the [EVM Deep Dive Part 6](https://noxx.substack.com/p/evm-deep-dives-the-path-to-shadowy-16e) article 