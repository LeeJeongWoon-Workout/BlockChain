package main

import (
	"fmt"
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
	"math/big"
	"math"
)

const targetBits=24

func IntToHex(n int64) []byte{
	return []byte(strconv.FormatInt(n,16))
}

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce		int64
}

type Blockchain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
			[][]byte{
					pow.block.PrevBlockHash,
					pow.block.Data,
					IntToHex(pow.block.Timestamp),
					IntToHex(int64(targetBits)),
					IntToHex(int64(nonce)),
			},
			[]byte{},
	)
	return data
}



func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64
	nonce=0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < math.MaxInt64 {
			data := pow.prepareData(nonce)
			hash = sha256.Sum256(data)
			fmt.Printf("\r%x", hash)

			hashInt.SetBytes(hash[:])
			if hashInt.Cmp(pow.target) == -1 {
					break
			} else {
					nonce++
			}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func NewracoonBlock() *Block {
	return NewBlock("Racoon coin", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewracoonBlock()}}
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("nonce: %d\n",block.Nonce)
			fmt.Println()
	}
}
