package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const targetBits = 16 // 难度值，越小越难

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 创建 PoW
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	return &ProofOfWork{b, target}
}

// 准备参与哈希的数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			HashTransactions(pow.block.Transactions),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// 挖矿
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("⛏️  开始挖矿...")

	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Printf("✅ 挖矿成功: %x\n", hash)
	return nonce, hash[:]
}

// 校验区块是否合法
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}
