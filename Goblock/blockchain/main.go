package main

import (
	"fmt"
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
	"math/big"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"log"
	"math"
)

const targetBits=24

func IntToHex(n int64) []byte{
	return []byte(strconv.FormatInt(n,16))
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block
}

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce		int64
}

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	i.currentHash = block.PrevBlockHash
	return block
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
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			lastHash = b.Get([]byte("l"))
			return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			err = b.Put([]byte("l"), newBlock.Hash)
			bc.tip = newBlock.Hash
			return nil
	})
}

func NewracoonBlock() *Block {
	return NewBlock("Racoon coin", []byte{})
}

const (
	blocksBucket="blocks"
	dbFile="chain.db"
)

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			if b == nil {
					genesis := NewracoonBlock()
					b, err := tx.CreateBucket([]byte(blocksBucket))
					err = b.Put(genesis.Hash, genesis.Serialize())
					err = b.Put([]byte("l"), genesis.Hash)
					tip = genesis.Hash
			} else {
					tip = b.Get([]byte("l"))
			}
			return nil
	})

	bc := Blockchain{tip, db}
	return &bc
}
