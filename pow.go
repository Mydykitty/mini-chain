package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

const targetBits = 16

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for {
		data := bytes.Join(
			[][]byte{
				pow.Block.PrevHash,
				pow.Block.HashTransactions(),
				IntToHex(pow.Block.Timestamp),
				IntToHex(int64(nonce)),
			}, []byte{},
		)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(pow.Block.Nonce)),
		}, []byte{},
	)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.Target) == -1
}
