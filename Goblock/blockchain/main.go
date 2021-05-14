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
const Difficulty=16

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n,16))
}
type Block struct {
	Data []byte
	Hash []byte
	PrevHash []byte
	nonce int64
	TS int64
}

type BlockChain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block)*ProofOfWork {
	target:=big.NewInt(1)
	target.Lsh(target,256-Difficulty)
	return &ProofOfWork{b,target}
}

func (pow *ProofOfWork) PrePareData(nonce int64) []byte {
	data:=bytes.Join([][]byte{
		pow.block.Data,
		pow.block.PrevHash,
		IntToHex(nonce),
		IntToHex(pow.block.TS),
		IntToHex(int64(Difficulty))},[]byte{})

	return data
}

func (pow *ProofOfWork)Run() (int64,[]byte) {
	var nonce int64
	var hashInt big.Int
	var hash [32]byte

	nonce=0

	for nonce<math.MaxInt64 {
		data:=pow.PrePareData(nonce)
		hash=sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target)==-1 {
			break
		}
		nonce++
	}
	return nonce,hash[:]
}

func CreateBlock(data string,prevhash []byte) *Block {
	block:=&Block{[]byte(data),[]byte{},prevhash,0,time.Hour.Microseconds()}
	pow:=NewProofOfWork(block)
	block.nonce,block.Hash=pow.Run()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prev:=chain.blocks[len(chain.blocks)-1]
	new:=CreateBlock(data,prev.Hash)
	chain.blocks=append(chain.blocks,new)
}

func Racoon() *Block{
	return CreateBlock("Racoon",[]byte{})
}

func InitBlockChain() *BlockChain{
	return &BlockChain{[]*Block{Racoon()}}
}

func main() {
	chain:=InitBlockChain()

	chain.AddBlock("Ivan 1 coin")
	chain.AddBlock("Fermat 3 coin")

	for _,block :=range chain.blocks {
		fmt.Printf("PrevHash: %x\n",block.PrevHash)
		fmt.Printf("Data: %s\n",block.Data)
		fmt.Printf("nonce: %d\n",block.nonce)

		fmt.Printf("\n")
	}
}