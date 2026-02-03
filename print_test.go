package main

import "testing"

func TestPrint(t *testing.T) {
	bc := CreateBlockchain()
	defer bc.DB.Close()

	bc.PrintChain()
}
