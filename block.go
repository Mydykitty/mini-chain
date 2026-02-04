package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
	Height       int // 👈 新增：区块高度
}

// 创建新区块
func NewBlock(txs []*Transaction, prevHash []byte, timeUnix int64, height int) *Block {
	block := &Block{
		Timestamp:    timeUnix,
		Transactions: txs,
		PrevHash:     prevHash,
		Hash:         []byte{},
		Nonce:        0,
		Height:       height,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// 创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 1700000000, 1)
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
