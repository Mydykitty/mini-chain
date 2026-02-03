package main

import "fmt"

func main() {
	bc := CreateBlockchain()
	defer bc.DB.Close()

	alice := NewWallet()
	bob := NewWallet()

	fmt.Printf("Alice地址: %x\n", alice.GetAddress())
	fmt.Printf("Bob地址: %x\n", bob.GetAddress())

	tx := NewTransaction(alice, bob.GetAddress(), 10)

	err := bc.MineBlock([]Transaction{*tx})
	if err != nil {
		fmt.Println("挖矿失败:", err)
		return
	}

	fmt.Println("挖矿成功")
}
