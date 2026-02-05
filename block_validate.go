package main

import "errors"

func ValidateBlock(block *Block, bc *Blockchain) error {
	// 1️⃣ 必须有交易
	if len(block.Transactions) == 0 {
		return errors.New("区块没有交易")
	}

	// 2️⃣ 第 0 笔必须是 coinbase
	if !block.Transactions[0].IsCoinbase() {
		return errors.New("coinbase 必须是第一笔交易")
	}

	// 3️⃣ 只能有 1 笔 coinbase
	coinbaseCount := 0
	for _, tx := range block.Transactions {
		if tx.IsCoinbase() {
			coinbaseCount++
		}
	}
	if coinbaseCount != 1 {
		return errors.New("区块 coinbase 数量不正确")
	}

	// 4️⃣ 校验 PoW
	pow := NewProofOfWork(block)
	if !pow.Validate() {
		return errors.New("PoW 校验失败")
	}

	// 5️⃣ 校验普通交易
	for i, tx := range block.Transactions {
		if i == 0 {
			continue
		}
		if !bc.VerifyTransaction(tx) {
			return errors.New("区块包含非法交易")
		}
	}

	return nil
}
