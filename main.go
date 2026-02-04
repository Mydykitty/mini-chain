package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	nodeID := os.Args[1]

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "mine":
			MineBlock(nodeID)
			return
		case "getchain":
			PrintChain(nodeID)
			return
		case "send":
			from := os.Args[3]
			to := os.Args[4]
			amount, _ := strconv.Atoi(os.Args[5])
			Send(from, to, amount, nodeID)
			return
		}
	}

	address := "miner"
	bc := CreateBlockchain(address, nodeID)
	go func() {
		fmt.Println("🟢 自动挖矿线程已启动")

		for {
			time.Sleep(10 * time.Second)

			if len(mempool) == 0 {
				continue
			}

			var txs []*Transaction

			for _, tx := range mempool {
				txs = append(txs, &tx)
			}

			cbTx := NewCoinbaseTX("miner-"+nodeID, "") // todo 需要加nodeID么
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)

			fmt.Println("⛏️ 打包交易挖出新区块，高度:", newBlock.Height)

			// 清空已打包交易
			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					SendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}
		}
	}()

	/*go func() {
		for {
			time.Sleep(20 * time.Second)
			MineAndBroadcastBlock(bc, nodeID)
		}
	}()*/
	StartServer(nodeID, bc)

}

/*
func main() {
	// 创建钱包
	aliceWallet := NewWallet()
	bobWallet := NewWallet()

	aliceAddr := string(aliceWallet.PublicKey)
	bobAddr := string(bobWallet.PublicKey)

	// 创建区块链
	bc := CreateBlockchain(aliceAddr)

	// 创世块余额
	fmt.Println("Alice创世块余额:", 100)

	// Alice转账给Bob 30
	tx := NewUTXOTransaction(aliceAddr, bobAddr, 30, bc)
	bc.AddBlock([]*Transaction{tx})

	fmt.Println("Alice -> Bob 30 转账完成")

	// 遍历区块链打印
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("\n=== 区块 ===\nHash: %x\nPrevHash: %x\n", block.Hash, block.PrevHash)
		for i, tx := range block.Transactions {
			fmt.Printf("  交易 %d\n", i)
			for j, out := range tx.Vout {
				fmt.Printf("    输出 %d 金额: %d\n", j, out.Value)
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	bc.PrintBlockchain()
}
*/
