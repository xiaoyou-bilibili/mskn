package proto

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
)

// Serialize proto序列化
func Serialize(data thrift.TStruct) []byte {
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocolConf(buf, nil)
	err := data.Write(context.Background(), proto)
	if err != nil {
		log.Printf("write data err %v", err)
		return nil
	}
	return buf.Bytes()
}

// Deserialize proto反序列化
func Deserialize(data []byte, res thrift.TStruct) {
	buf := thrift.NewTMemoryBuffer()
	_, err := buf.Write(data)
	if err != nil {
		log.Printf("write data err %v", err)
		return
	}
	proto := thrift.NewTBinaryProtocolConf(buf, nil)
	err = res.Read(context.Background(), proto)
	if err != nil {
		log.Printf("read data err %v", err)
		return
	}
}
