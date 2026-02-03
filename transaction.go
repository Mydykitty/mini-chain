package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"math/big"
)

type Transaction struct {
	From      []byte
	To        []byte
	Amount    int
	Signature []byte
	PubKey    []byte
}

func NewTransaction(w *Wallet, to []byte, amount int) *Transaction {
	tx := &Transaction{
		From:   w.GetAddress(),
		To:     to,
		Amount: amount,
		PubKey: w.PublicKey,
	}
	tx.Sign(w.PrivateKey)
	return tx
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.Signature = nil
	txCopy.PubKey = nil

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(txCopy)
	hash = sha256.Sum256(buf.Bytes())
	return hash[:]
}

func (tx *Transaction) Sign(priv ecdsa.PrivateKey) {
	hash := tx.Hash()
	r, s, _ := ecdsa.Sign(rand.Reader, &priv, hash)
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature
}

func (tx *Transaction) Verify() bool {
	hash := tx.Hash()

	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:sigLen/2])
	s.SetBytes(tx.Signature[sigLen/2:])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(tx.PubKey)
	x.SetBytes(tx.PubKey[:keyLen/2])
	y.SetBytes(tx.PubKey[keyLen/2:])

	pubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
	return ecdsa.Verify(&pubKey, hash, &r, &s)
}

func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(tx)
	return buf.Bytes()
}
