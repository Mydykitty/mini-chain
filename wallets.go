package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet // key: 地址，value: 钱包对象
}

const walletFile = "wallets.json"

// 创建钱包管理器
func NewWallets() *Wallets {
	wallets := &Wallets{Wallets: make(map[string]*Wallet)}

	if _, err := os.Stat(walletFile); err == nil {
		data, err := ioutil.ReadFile(walletFile)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal(data, wallets)
		if err != nil {
			log.Panic(err)
		}
	}

	return wallets
}

// 保存钱包到文件
func (ws *Wallets) Save() {
	data, err := json.Marshal(ws)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(walletFile, data, 0644)
}

// 创建一个新钱包
func (ws *Wallets) CreateWallet(name string) string {
	wallet := NewWallet()
	ws.Wallets[name] = wallet
	ws.Save()
	fmt.Println("✅ 钱包创建成功，名称:", name)
	return name
}

// 获取钱包
func (ws *Wallets) GetWallet(name string) *Wallet {
	return ws.Wallets[name]
}

// 列出所有钱包
func (ws *Wallets) ListWallets() {
	fmt.Println("==== 钱包列表 ====")
	for name := range ws.Wallets {
		fmt.Println(name)
	}
}
