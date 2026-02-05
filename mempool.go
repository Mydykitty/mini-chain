package main

import (
	"encoding/hex"
	"fmt"
	"sync"
)

type Mempool struct {
	mu  sync.RWMutex
	txs map[string]Transaction
}

var mempool = Mempool{
	txs: make(map[string]Transaction),
}

func (mp *Mempool) Remove(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	delete(mp.txs, txID)
}

func (mp *Mempool) GetTransactions() []*Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	var txs []*Transaction
	for _, tx := range mp.txs {
		txCopy := tx
		txs = append(txs, &txCopy)
	}
	return txs
}

func (mp *Mempool) RemoveInvalid(bc *Blockchain) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for id, tx := range mp.txs {
		if !bc.VerifyTransaction(&tx) {
			fmt.Println("🗑 移除失效交易:", id)
			delete(mp.txs, id)
		}
	}
}

func (mp *Mempool) AddToMempool(tx Transaction, bc *Blockchain) {
	if bc.VerifyTransaction(&tx) {
		mp.txs[hex.EncodeToString(tx.ID)] = tx
	} else {
		fmt.Println("❌ 非法交易，拒绝加入mempool")
	}
}
