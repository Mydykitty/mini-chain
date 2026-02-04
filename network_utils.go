package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func commandToBytes(command string) []byte {
	var bytes [12]byte
	copy(bytes[:], command)
	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("❌ 无法连接节点 %s\n", addr)
		return
	}
	defer conn.Close()
	conn.Write(data)
}

func SendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}
