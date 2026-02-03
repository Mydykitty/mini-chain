package main

import "encoding/binary"

func IntToHex(num int64) []byte {
	buff := make([]byte, 8)
	binary.BigEndian.PutUint64(buff, uint64(num))
	return buff
}
