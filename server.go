package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const protocol = "tcp"
const nodeVersion = 1

var nodeAddress string
var knownNodes = []string{"localhost:3000"} // 种子节点
var miningInterrupt = make(chan bool)

func StartServer(nodeID string, bc *Blockchain) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	fmt.Printf("🌐 节点已启动: %s\n", nodeAddress)

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *Blockchain) {
	defer conn.Close()
	request, err := io.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:12])
	fmt.Printf("📩 收到命令: %s\n", command)

	switch command {
	case "version":
		handleVersion(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "block":
		handleBlock(request, bc)
	case "tx":
		handleTx(request, bc)
	default:
		fmt.Println("未知命令")
	}
}
