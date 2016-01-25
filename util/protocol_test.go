package util

import (
	//"fmt"
	//m_util "github.com/alexwang58/meipush/util"
	//"bytes"
	"testing"
)

func TestPackAndDepack(t *testing.T) {
	message := "Hello World"
	pack := AssemblePack([]byte(message))
	t.Logf("[HEAD] : %b", []byte(HEAD_PACK))
	t.Logf("[HEAD_Lenght] : %d", len(HEAD_PACK))
	if len(HEAD_PACK) != HEAD_PACK_LEN {
		t.Fatalf("HEAD_PACK_LEN want %d , given %d\n", len(HEAD_PACK), HEAD_PACK_LEN)
	}
	t.Logf("[Message_Lenght] : %d", len(message))
	t.Logf("[Message_Lenght_Byte] : %b", int2Byte(len(message)))
	t.Logf("[PACK_Byte] : %b", pack)
	t.Log(HEAD_PACK_BYTES)

	readerChannel := make(chan []byte, 16)
	tmp := DisassemblePack(pack, readerChannel)
	data := <-readerChannel
	if string(data) != message {
		t.Fatalf("DissemblePack Error \n[want] %s\n[given] %s", message, string(data))
	}
	t.Log(tmp, data)
}
