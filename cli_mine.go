package main

import "fmt"

func MineBlock(nodeID string) {
	bc := NewBlockchain(nodeID)
	defer bc.DB.Close()

	cbTx := NewCoinbaseTX("miner", "CLI Mine Reward")
	block := bc.MineBlock([]*Transaction{cbTx})

	fmt.Println("⛏️  挖出新区块成功!")
	fmt.Printf("高度: %d\n", block.Height)
	fmt.Printf("Hash: %x\n", block.Hash)
}

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

func MineAndBroadcastBlock(bc *Blockchain, nodeID string) {
	fmt.Println("⛏️  节点", nodeID, "开始挖矿...")

	cbTx := NewCoinbaseTX("miner", "")
	newBlock := bc.MineBlock([]*Transaction{cbTx})

	fmt.Println("✅ 挖矿成功，高度:", newBlock.Height)

	for _, node := range knownNodes {
		if node != "localhost:"+nodeID {
			SendInv(node, "block", [][]byte{newBlock.Hash})
		}
	}
}
