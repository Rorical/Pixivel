package config

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

//Hash will md5 the struct
func HashStruct(item interface{}) string {
	jsonBytes, _ := json.Marshal(item)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

func StringOut(bye []byte) string {
	return string(bye)
}

func StringIn(strings string) []byte {
	return []byte(strings)
}

func IntIn(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func IntOut(bye []byte) int {
	bytebuff := bytes.NewBuffer(bye)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func UintIn(n uint) []byte {
	data := uint64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func UintOut(bye []byte) uint {
	bytebuff := bytes.NewBuffer(bye)
	var data uint64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return uint(data)
}
