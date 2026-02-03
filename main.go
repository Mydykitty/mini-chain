package main

import "fmt"

func main() {
	bc := CreateBlockchain()
	defer bc.DB.Close()

	fmt.Println("Alice 余额:", bc.GetBalance("Alice"))
	fmt.Println("Bob 余额:", bc.GetBalance("Bob"))

	tx1 := Transaction{"Alice", "Bob", 50}
	tx2 := Transaction{"Bob", "Charlie", 20}

	err := bc.MineBlock([]Transaction{tx1, tx2})
	if err != nil {
		fmt.Println("❌ 挖矿失败:", err)
		return
	}
	fmt.Println("✅ 区块挖矿成功")

	fmt.Println("挖矿后余额：")
	fmt.Println("Alice:", bc.GetBalance("Alice"))
	fmt.Println("Bob:", bc.GetBalance("Bob"))
	fmt.Println("Charlie:", bc.GetBalance("Charlie"))
}
