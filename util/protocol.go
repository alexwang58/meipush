package util

import (
	"bytes"
	"encoding/binary"
	//"fmt"
)

const (
	HEAD_PACK        = "MEIMSG/1.0\r\n"
	MSG_LEN_METADATA = 4 // 消息长度用4个byte表示
	HEAD_PACK_LEN    = len(HEAD_PACK)
	//HEAD_PACK_BYTES  = []byte(HEAD_PACK)
)

var HEAD_PACK_BYTES = []byte(HEAD_PACK)

//HEAD_PACK_BYTES = []byte("dddd")

func AssemblePack(message []byte) []byte {
	return append(append(HEAD_PACK_BYTES, int2Byte(len(message))...), message...)
}

// 从传入的buffer中解析message，并写入readerChannel
func DisassemblePack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i++ {
		if length < i+HEAD_PACK_LEN+MSG_LEN_METADATA {
			//fmt.Println(i, length, "Break")
			// 输入长度<协议头长度
			break
		}
		// 刚好为消息头
		if bytes.Equal(HEAD_PACK_BYTES, buffer[i:i+HEAD_PACK_LEN]) {
			metaLen := i + HEAD_PACK_LEN + MSG_LEN_METADATA
			messageLength := byte2Int(buffer[i+HEAD_PACK_LEN : metaLen])
			totalLen := i + HEAD_PACK_LEN + MSG_LEN_METADATA + messageLength
			if length < totalLen {
				// 输入长度<协议头描述的消息长度
				//fmt.Println("Break")
				break
			}
			message := buffer[metaLen:totalLen]
			readerChannel <- message
		}
	}

	if i == length {
		return make([]byte, 0)
	}

	return buffer[i:] // 将不完整的消息返回
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
