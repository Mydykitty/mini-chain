package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

func CreateBlockchain() *Blockchain {
	db, _ := bolt.Open(dbFile, 0600, nil)

	var tip []byte
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, _ = tx.CreateBucket([]byte(blocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}

func (bc *Blockchain) MineBlock(txs []Transaction) error {
	for _, tx := range txs {
		if !tx.Verify() {
			return fmt.Errorf("发现非法签名交易")
		}
		/*if len(tx.From) == 0 { // 挖矿奖励
			continue
		}
		if bc.GetBalance(tx.From) < tx.Amount {
			return fmt.Errorf("交易非法，余额不足: %s", tx.From)
		}*/
	}

	var lastHash []byte
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	block := NewBlock(txs, lastHash)

	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(block.Hash, block.Serialize())
		b.Put([]byte("l"), block.Hash)
		bc.Tip = block.Hash
		return nil
	})

	return nil
}

func NewBlock(txs []Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), txs, prevHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	genesisTx := Transaction{
		From:      nil,
		To:        []byte("genesis"), // 只是占位标识
		Amount:    100,
		Signature: nil,
		PubKey:    nil,
	}
	return NewBlock([]Transaction{genesisTx}, []byte{})
}

func (bc *Blockchain) GetBalance(address []byte) int {
	balance := 0

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			block := DeserializeBlock(v)

			for _, tx := range block.Transactions {
				if tx.To != nil && bytes.Equal(tx.To, address) {
					balance += tx.Amount
				}
				if tx.From != nil && bytes.Equal(tx.From, address) {
					balance -= tx.Amount
				}
			}
		}
		return nil
	})

	return balance
}
