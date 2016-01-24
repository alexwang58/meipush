package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	HEAD_PACK     = "MeIhEaDeR/"
	HEAD_PACK_LEN = len(HEAD_PACK)
)

func AssemblePack(message []byte) []byte {
	return append(append([]byte(HEAD_PACK), int2Byte(len(message))...), message...)
}

func DisassemblePack(buffer []byte) {
	length := len(buffer)
	fmt.Println(length)
}

//打印形式为IP表达形式[255, 255, 255, 255]
func int2Byte(n int) []byte {
	x := int32(n)
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, x)
	return byteBuffer.Bytes()
}

func byte2Int(b []byte) int {
	var x int32
	byteBuffer := bytes.NewBuffer(b)
	binary.Read(byteBuffer, binary.BigEndian, &x)
	return int(x)
}
