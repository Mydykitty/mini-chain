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

			fmt.Printf("\n====== 区块 ======\n")
			fmt.Printf("时间: %d\n", block.Timestamp)
			fmt.Printf("PrevHash: %x\n", block.PrevHash)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("Nonce: %d\n", block.Nonce)

			for i, tx := range block.Transactions {
				fmt.Printf("  --- 交易 %d ---\n", i)

				if tx.From == nil {
					fmt.Printf("    From: <系统奖励>\n")
				} else {
					fmt.Printf("    From: %x\n", tx.From)
				}

				fmt.Printf("    To:   %x\n", tx.To)
				fmt.Printf("    金额: %d\n", tx.Amount)

				if tx.Signature != nil {
					valid := tx.Verify()
					fmt.Printf("    签名有效: %v\n", valid)
				} else {
					fmt.Printf("    签名: <无>\n")
				}
			}
		}
		return nil
	})
}
