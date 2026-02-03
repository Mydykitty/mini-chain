package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return &Wallet{*priv, pub}
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	ripemd := ripemd160.New()
	ripemd.Write(publicSHA256[:])
	return ripemd.Sum(nil)
}

func (w *Wallet) GetAddress() []byte {
	pubHash := HashPubKey(w.PublicKey)
	return pubHash
}
