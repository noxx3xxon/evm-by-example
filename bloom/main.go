package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {

	// Set up eth client
	client, err := ethclient.Dial(os.Getenv("ETH_RPC_URL"))
	checkError(err)

	// Get block data
	blockNumber := big.NewInt(15001871)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	checkError(err)

	blockHeaderJSON, err := json.MarshalIndent(block.Header(), "", "  ")
	fmt.Printf("Block Header for Block %d \n\n%s \n\n", blockNumber, string(blockHeaderJSON))

	transactionJSON, err := json.MarshalIndent(block.Transactions()[0], "", "  ")
	fmt.Printf("Transaction Data for Tx %s \n\n%s \n\n", block.Transactions()[0].Hash().Hex(), string(transactionJSON))

	blockReceipts := []*types.Receipt{}

	// Using tx's in block fill the Receipts slice
	for _, tx := range block.Transactions() {
		txReceipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		checkError(err)
		blockReceipts = append(blockReceipts, txReceipt)

	}

	transactionReceiptJSON, err := json.MarshalIndent(blockReceipts[0], "", "  ")
	checkError(err)

	transactionReceiptLogsJSON, err := json.MarshalIndent(blockReceipts[0].Logs, "", "  ")
	checkError(err)

	fmt.Printf("Transaction Receipt for Tx %s \n\n%s \n\n", blockReceipts[0].TxHash.Hex(), string(transactionReceiptJSON))

	fmt.Printf("Event Logs for Tx %s \n\n%s \n\n", blockReceipts[0].TxHash.Hex(), string(transactionReceiptLogsJSON))

	// Get Byte & Bit Indexes that need to be updated in the bloom filter for tx hash 0x311ba3a0affb00510ae3f0a36c5bcd0a48cdb23d803bbc16f128639ffb9e3e58
	for _, log := range blockReceipts[0].Logs {
		buf := make([]byte, 6)

		// Get Byte & Bit indexes for log address
		_, _, _, _, _, _ = bloomValues(log.Address.Bytes(), buf)

		// Get Byte & Bit indexes for each topic
		for _, topic := range log.Topics {
			_, _, _, _, _, _ = bloomValues(topic[:], buf)
		}

	}

	// Create a bloom filter using the public method available in Geth
	receiptBloom := types.CreateBloom(blockReceipts[:1])

	// Print the hex and binary version of the bloom filter
	fmt.Printf("Transaction Receipt logsBloom (hex)\n %s \n\n", hex.EncodeToString(receiptBloom.Bytes()))
	fmt.Printf("Transaction Receipt logsBloom (binary)\n %s \n\n", fmt.Sprintf("%08b", receiptBloom.Bytes()))

	// Printing the bytes associated with the Event signature topic for cross reference
	fmt.Printf("Transaction Receipt logsBloom byte index 75\n %s \n", fmt.Sprintf("%08b", receiptBloom.Bytes()[75]))
	fmt.Printf("Transaction Receipt logsBloom byte index 195\n %s \n", fmt.Sprintf("%08b", receiptBloom.Bytes()[195]))
	fmt.Printf("Transaction Receipt logsBloom byte index 123\n %s \n", fmt.Sprintf("%08b", receiptBloom.Bytes()[123]))
}

// hasherPool holds LegacyKeccak256 hashers for rlpHash.
var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// bloomValues is a private function in geth I've copied it and added in some print statements to show you how it works
// See here for function in geth codebase https://github.com/ethereum/go-ethereum/blob/d8ff53dfb8a516f47db37dbc7fd7ad18a1e8a125/core/types/bloom9.go#L139
func bloomValues(data []byte, hashbuf []byte) (uint, byte, uint, byte, uint, byte) {
	sha := hasherPool.Get().(crypto.KeccakState)
	sha.Reset()
	sha.Write(data)
	sha.Read(hashbuf)
	hasherPool.Put(sha)

	fmt.Printf("data []byte: %s \n", hex.EncodeToString(data))
	fmt.Printf("hashbuf []byte (Post hashed data being loaded in): %s \n", hex.EncodeToString(hashbuf))

	// The actual bits to flip
	v1 := byte(1 << (hashbuf[1] & 0x7))
	v2 := byte(1 << (hashbuf[3] & 0x7))
	v3 := byte(1 << (hashbuf[5] & 0x7))

	// The indices for the bytes to OR in
	i1 := types.BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3) - 1
	i2 := types.BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3) - 1
	i3 := types.BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3) - 1

	fmt.Printf(
		`
BYTE PAIR 1 - (%s)
	BIT INDEX / BYTE VALUE
		%s = %s = hashbuf[1]
		%s & %s = %s = %d
		%s << %d = %s = v1
	BYTE INDEX
		%s = %s
		%s & %s = %s
		%s >> 3 = %s = %d
		%d - %d - 1 = %d = i1

`,
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf)),
		fmt.Sprintf("%x", hashbuf[1]), fmt.Sprintf("%08b", hashbuf[1]),
		fmt.Sprintf("%08b", hashbuf[1]), fmt.Sprintf("%08b", 0x7), fmt.Sprintf("%08b", (hashbuf[1]&0x7)), (hashbuf[1] & 0x7),
		fmt.Sprintf("%08b", (hashbuf[1]&0x7)), (hashbuf[1] & 0x7), fmt.Sprintf("%08b", v1),
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf)), fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf)),
		fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf)), fmt.Sprintf("%08b", 0x7ff), fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf)&0x7ff)),
		fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf)&0x7ff)), fmt.Sprintf("%08b", ((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3)), uint((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3),
		types.BloomByteLength, uint((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3), i1,
	)

	fmt.Printf(
		`
BYTE PAIR 2 - (%s)
	BIT INDEX / BYTE VALUE
		%s = %s = hashbuf[3]
		%s & %s = %s = %d
		%s << %d = %s = v2
	BYTE INDEX
		%s = %s
		%s & %s = %s
		%s >> 3 = %s = %d
		%d - %d - 1 = %d = i2

`,
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf[2:])),
		fmt.Sprintf("%x", hashbuf[3]), fmt.Sprintf("%08b", hashbuf[3]),
		fmt.Sprintf("%08b", hashbuf[3]), fmt.Sprintf("%08b", 0x7), fmt.Sprintf("%08b", (hashbuf[3]&0x7)), (hashbuf[3] & 0x7),
		fmt.Sprintf("%08b", (hashbuf[3]&0x7)), (hashbuf[3] & 0x7), fmt.Sprintf("%08b", v1),
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf[2:])), fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf[2:])),
		fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf[2:])), fmt.Sprintf("%08b", 0x7ff), fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)),
		fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)), fmt.Sprintf("%08b", ((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3)), uint((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3),
		types.BloomByteLength, uint((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3), i1,
	)

	fmt.Printf(
		`
BYTE PAIR 3 - (%s)
	BIT INDEX / BYTE VALUE
		%s = %s = hashbuf[5]
		%s & %s = %s = %d
		%s << %d = %s = v3
	BYTE INDEX
		%s = %s
		%s & %s = %s
		%s >> 3 = %s = %d
		%d - %d - 1 = %d = i3

`,
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf[4:])),
		fmt.Sprintf("%x", hashbuf[5]), fmt.Sprintf("%08b", hashbuf[5]),
		fmt.Sprintf("%08b", hashbuf[5]), fmt.Sprintf("%08b", 0x7), fmt.Sprintf("%08b", (hashbuf[5]&0x7)), (hashbuf[5] & 0x7),
		fmt.Sprintf("%08b", (hashbuf[5]&0x7)), (hashbuf[5] & 0x7), fmt.Sprintf("%08b", v1),
		fmt.Sprintf("%x", binary.BigEndian.Uint16(hashbuf[4:])), fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf[4:])),
		fmt.Sprintf("%08b", binary.BigEndian.Uint16(hashbuf[4:])), fmt.Sprintf("%08b", 0x7ff), fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)),
		fmt.Sprintf("%08b", (binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)), fmt.Sprintf("%08b", ((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3)), uint((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3),
		types.BloomByteLength, uint((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3), i1,
	)

	return i1, v1, i2, v2, i3, v3
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
