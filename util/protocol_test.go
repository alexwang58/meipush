package util

import (
	//"fmt"
	//m_util "github.com/alexwang58/meipush/util"
	"testing"
)

func TestPackAndDepack(t *testing.T) {
	//fmt.Println(m_util.AssemblePack([]byte("Hello World")))
	message := "Readreadsstructuredbinarydatafromrintodata.Datamustbeapointertoafixed-sizevalueorasliceoffixed-sizevalues.Bytesreadfromraredecodedusingthespecifiedbyteorderandwrittentosuccessivefieldsofthedata.Whenreadingintostructs,thefielddataforfieldswithblankfieldnamesisskipped;i.e.,blankfieldnames"
	pack := AssemblePack([]byte(message))
	t.Logf("[HEAD] : %b", []byte(HEAD_PACK))
	t.Logf("[HEAD_Lenght] : %d", len(HEAD_PACK))
	if len(HEAD_PACK) != HEAD_PACK_LEN {
		t.Fatalf("HEAD_PACK_LEN want %d , given %d\n", len(HEAD_PACK), HEAD_PACK_LEN)
	}
	t.Logf("[Message_Lenght] : %d", len(message))
	t.Logf("[Message_Lenght_Byte] : %b", int2Byte(len(message)))
	t.Logf("[PACK_Byte] : %b", pack)
	//t.Logf("[MESSAGE_LEN] : %b", m_util.int2Byte(len(message)))
	// 0 0 0 5
	//byt := [...]byte{77, 101, 73, 104, 69, 97, 68, 101, 82, 47, 0, 0, 0, 11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}
	//t.Log(byt)
	/*
		if byt == pack {
			t.Log("Same")
		}
		m_
	*/
	DisassemblePack(pack)
	//t.Errorf("util.AssemblePack() encode package error:  %s | GIVEN %s", pack, okencode)
}
