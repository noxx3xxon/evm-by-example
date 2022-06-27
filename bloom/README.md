# BLOOM FILTERS

This section covers the bloom filters, a key data structure used to create the logsBloom in both the block header and transaction receipts.

See the the "Bloom Filter" section in this [EVM Deep Dive](https://noxx.substack.com/p/evm-deep-dives-the-path-to-shadowy-16e?utm_source=%2Fprofile%2F80455042-noxx&utm_medium=reader2) article for more details.

To run the code associated with this section first ensure the ETH_RPC_URL is set then run

```
go run bloom/main.go
```

## Example - Block 15001871

This code section run through the bloom filters contained in Block 15001871. The following print statements have been added to take you from block header to the specific bits that are flipped within the bloom filter.

- Prints Block Header for 15001871
- Prints Transaction Data for first Tx in Block 15001871 0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58
- Prints Transaction Receipt for 0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58
- Prints Event Logs for Transaction Receipt 0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58
- Prints out bit flip indexes in bloom filter caused by the Event Log of Transaction Receipt 0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58
- Prints out Bloom Filter in hexadecimal
- Prints out Bloom filter in binary
- Prints out specific bytes within Bloom Filter that were referenced in the [EVM Deep Dive](https://noxx.substack.com/p/evm-deep-dives-the-path-to-shadowy-16e?utm_source=%2Fprofile%2F80455042-noxx&utm_medium=reader2) article 