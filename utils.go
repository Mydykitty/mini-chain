package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// Sha256 计算哈希
func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HashPubKey 对公钥做 HASH160（SHA256 + RIPEMD160）
func HashPubKey(pubKey []byte) []byte {
	pubSHA256 := Sha256(pubKey)

	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(pubSHA256)
	pubRIPEMD160 := ripemdHasher.Sum(nil)

	return pubRIPEMD160
}

// int64 转 []byte
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// 把区块中所有交易做一次hash
func HashTransactions(txs []*Transaction) []byte {
	var txHashes [][]byte
	for _, tx := range txs {
		txHashes = append(txHashes, tx.ID)
	}
	data := bytes.Join(txHashes, []byte{})
	hash := Sha256(data)
	return hash
}
