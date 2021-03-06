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
func IntToHex(n int64) []byte{
	return []byte(strconv.FormatInt(n,16))
}
const targetBits = 24

type ProofOfWork struct {
	block *Block
	target *big.Int
}



type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}
type BlockChain struct {
	blocks []*Block
}
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
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

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce <	math.MaxInt64 {
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

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func main() {
	chain:=NewBlockchain()
	chain.AddBlock("1 coin")
	chain.AddBlock("2 coin")
	chain.AddBlock("3 coin")
	for _, block := range chain.blocks {

			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
	}
}
