package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func decodePayload(data []byte, v interface{}) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(v)
}

func sendVersion(addr string, bc *Blockchain) {
	bestHeight := getBestHeight(bc)
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

func handleVersion(request []byte, bc *Blockchain) {
	var payload Version
	decodePayload(request[12:], &payload)

	myHeight := getBestHeight(bc)

	if payload.BestHeight > myHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if payload.BestHeight < myHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

func SendInv(addr, kind string, items [][]byte) {
	payload := gobEncode(Inv{nodeAddress, kind, items})
	request := append(commandToBytes("inv"), payload...)
	sendData(addr, request)
}

func handleInv(request []byte, bc *Blockchain) {
	var payload Inv
	decodePayload(request[12:], &payload)

	if payload.Type == "block" {
		blocksInTransit = payload.Items
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
	}
}

func sendGetData(addr, kind string, id []byte) {
	payload := gobEncode(GetData{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)
	sendData(addr, request)
}

func handleGetData(request []byte, bc *Blockchain) {
	var payload GetData
	decodePayload(request[12:], &payload)

	if payload.Type == "block" {
		block, _ := bc.GetBlock(payload.ID)
		data := gobEncode(BlockData{nodeAddress, block.Serialize()})
		request := append(commandToBytes("block"), data...)
		sendData(payload.AddrFrom, request)
	}
}

func handleBlock(request []byte, bc *Blockchain) {
	var payload BlockData
	decodePayload(request[12:], &payload)

	block := DeserializeBlock(payload.Block)
	fmt.Println("⛓️ 收到新区块:", block.Hash)

	bc.AddBlockFromNetwork(block)

	fmt.Printf("当前区块高度: %d\n", getBestHeight(bc))
}
