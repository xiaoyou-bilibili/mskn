package proto

import (
	"fmt"
	"mskn-server/core/proto/mn"
	"mskn-server/core/proto/mnt"
	"testing"
)

func TestProto(t *testing.T) {
	code := mnt.CodePush{
		Addr: "123",
		Name: "456",
		Code: "789",
	}

	tmp := Serialize(&code)
	code2 := mnt.CodePush{}
	fmt.Println(code2)
	Deserialize(tmp, &code2)
	fmt.Println(code2)
}

func TestEnum(t *testing.T) {
	fmt.Println(mn.MessageTypeTaskAck)
}
