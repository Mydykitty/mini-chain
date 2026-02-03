package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func (bc *Blockchain) PrintChain() {
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			block := DeserializeBlock(v)

			fmt.Printf("==== 区块 ====\n")
			fmt.Printf("时间: %d\n", block.Timestamp)
			fmt.Printf("PrevHash: %x\n", block.PrevHash)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("Nonce: %d\n", block.Nonce)

			for _, tx := range block.Transactions {
				fmt.Printf("  交易: %s -> %s (%d)\n", tx.From, tx.To, tx.Amount)
			}
			fmt.Println()
		}
		return nil
	})
}
