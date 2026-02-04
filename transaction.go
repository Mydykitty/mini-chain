package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	OutIndex  int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

// Coinbase判断
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].OutIndex == -1
}

// 设置ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	tx.ID = Sha256(encoded.Bytes())
}

// 创建UTXO输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// 输出锁定到某个地址
func (out *TXOutput) Lock(address []byte) {
	out.PubKeyHash = HashPubKey(address)
}

// 判断是否属于某人
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// Coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = "Block reward"
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(100, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()
	return &tx
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}
