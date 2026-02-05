package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"math/big"
)

type Wallet struct {
	PrivateKeyBytes []byte
	PublicKey       []byte
}

func NewWallet() *Wallet {
	priv, pub := newKeyPair()
	privBytes := priv.D.Bytes()
	return &Wallet{
		PrivateKeyBytes: privBytes,
		PublicKey:       pub,
	}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return *priv, pub
}

func (w *Wallet) ToECDSA() *ecdsa.PrivateKey {
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(w.PrivateKeyBytes)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(w.PrivateKeyBytes)
	return priv
}

func (w *Wallet) Address() []byte {
	return HashPubKey(w.PublicKey)
}
