package main

import (
	"bytes"
	"encoding/gob"
)

type Transaction struct {
	From   string
	To     string
	Amount int
}

func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(tx)
	return buf.Bytes()
}
