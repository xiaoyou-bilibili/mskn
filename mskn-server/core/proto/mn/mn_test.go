package mn

import (
	"fmt"
	"testing"
)

func TestMessageType(t *testing.T) {
	mn := NewMn(EncryptTypeBase64)
	mn.setMessageType(MessageTypeConnect)
	fmt.Printf("%b", mn.msgType)
	mn.setMessageLen()
	fmt.Println(mn.msgLen)
	fmt.Println(mn.getSha256("123"))
}
