package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
}

// 创建新区块
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), txs, prevHash, []byte{}, 0}
	hash := block.HashBlock()
	block.Hash = hash
	return block
}

// 创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 序列化
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

// 反序列化
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 简单 hash 作为演示
func (b *Block) HashBlock() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	hash := Sha256(encoded.Bytes())
	return hash
}
