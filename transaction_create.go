package main

import (
	"encoding/hex"
	"log"
)

// 创建一笔UTXO转账交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	ws := NewWallets()
	// 确保钱包存在
	if ws.GetWallet(from) == nil {
		ws.CreateWallet(from)
	}
	if ws.GetWallet(to) == nil {
		ws.CreateWallet(to)
	}

	fromWallet := ws.GetWallet(from)
	toWallet := ws.GetWallet(to)

	if fromWallet == nil || toWallet == nil {
		log.Panic("钱包不存在")
	}

	fromPubKeyHash := fromWallet.Address()

	acc, validOutputs := bc.FindSpendableOutputs(fromPubKeyHash, amount)
	if acc < amount {
		log.Panic("余额不足")
	}

	var inputs []TXInput
	var outputs []TXOutput

	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)
		for _, outIdx := range outs {
			inputs = append(inputs, TXInput{
				Txid:     txID,
				OutIndex: outIdx,
				PubKey:   fromWallet.PublicKey,
			})
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, toWallet.Address()))

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, fromWallet.Address()))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	prevTXs := make(map[string]Transaction)
	for txid := range validOutputs {
		id, _ := hex.DecodeString(txid)
		prevTXs[txid], _ = bc.FindTransaction(id)
	}

	tx.Sign(fromWallet.ToECDSA(), prevTXs)
	return &tx
}
