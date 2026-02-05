package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	nodeID := os.Args[1]

	ws := NewWallets()
	if ws.GetWallet("miner") == nil {
		ws.CreateWallet("miner")
	}
	myWallet := ws.GetWallet("miner") // 或 miner-3000
	fmt.Println("address is: ", myWallet.Address())

	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "getchain":
			PrintChain(nodeID)
			return
		case "send":
			from := os.Args[3]
			to := os.Args[4]
			amount, _ := strconv.Atoi(os.Args[5])
			Send(from, to, amount, nodeID)
			return
		case "createwallet":
			name := os.Args[3]
			ws.CreateWallet(name)
			return
		case "listwallets":
			ws.ListWallets()
			return
		case "balance":
			name := os.Args[3]
			GetBalance(name, nodeID)
			return
		}
	}

	bc := CreateBlockchain(myWallet.Address(), nodeID)
	go func() {
		fmt.Println("🟢 自动挖矿线程已启动")

		for {
			time.Sleep(10 * time.Second)

			if len(mempool.GetTransactions()) == 0 {
				continue
			}

			cbTx := NewCoinbaseTX(myWallet.Address(), "")
			txs := []*Transaction{cbTx}

			for _, tx := range mempool.GetTransactions() {
				txs = append(txs, tx)
			}

			newBlock := bc.MineBlock(txs)
			if newBlock == nil {
				fmt.Println("⚠️ main本轮挖矿被中断")
				continue
			}

			fmt.Println("⛏️ 打包交易挖出新区块，高度:", newBlock.Height)

			// 清空已打包交易
			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				mempool.Remove(txID)
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
	StartTxServer(bc, nodeID)
	StartServer(nodeID, bc)

}

func StartTxServer(bc *Blockchain, port string) {
	http.HandleFunc("/tx", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Amount int    `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		tx := NewUTXOTransaction(req.From, req.To, req.Amount, bc)
		txID := hex.EncodeToString(tx.ID)

		SendTx("localhost:"+port, tx)

		fmt.Println("✅ 交易已创建并广播啊")

		fmt.Fprintf(w, "Transaction added: %s\n", txID)
	})

	go func() {
		fmt.Println("🌐 HTTP Server started at port", port)
		if err := http.ListenAndServe(":"+"4"+port, nil); err != nil { // 43000、43001 等端口
			log.Panic(err)
		}
	}()
}
