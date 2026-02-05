package main

import "fmt"

func PrintChain(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.DB.Close()

	bc.PrintBlockchain()
}

func Send(from, to string, amount int, nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.DB.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)

	SendTx("localhost:"+nodeID, tx)

	fmt.Println("✅ 交易已创建并广播")
}
