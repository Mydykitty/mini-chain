package main

import (
	"fmt"
)

func GetBalance(address string, nodeID string) int {
	bc := NewBlockchain(nodeID)
	defer bc.DB.Close()

	pubKeyHash := HashPubKey([]byte(address))

	// 查找所有 UTXO，amount 用大数保证查出全部
	balance, _ := bc.FindSpendableOutputs(pubKeyHash, 1<<31)

	fmt.Printf("💰 %s 余额: %d\n", address, balance)
	return balance
}
