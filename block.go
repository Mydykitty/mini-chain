package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Timestamp    int64
	Transactions []Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	gob.NewEncoder(&res).Encode(b)
	return res.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block
	gob.NewDecoder(bytes.NewReader(data)).Decode(&block)
	return &block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	hash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return hash[:]
}
