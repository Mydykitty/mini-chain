package main

import (
	"encoding/hex"
	"log"
)

// 创建一笔UTXO转账交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	// 获取钱包
	wallet := NewWallet() // 注意：在完整项目中你应该用钱包管理器获取已有钱包
	pubKeyHash := HashPubKey([]byte(from))

	// 查找可用 UTXO
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		log.Panic("余额不足")
	}

	// 构建输入
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, outIdx := range outs {
			input := TXInput{txID, outIdx, nil, []byte(from)}
			inputs = append(inputs, input)
		}
	}

	// 构建输出
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		// 找零给自己
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	// 创建交易
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	// 签名交易
	prevTXs := make(map[string]Transaction)
	for txid := range validOutputs {
		txInBytes, _ := hex.DecodeString(txid)
		prevTx, _ := bc.FindTransaction(txInBytes)
		prevTXs[txid] = prevTx
	}
	tx.Sign(wallet.PrivateKey, prevTXs)

	return &tx
}
