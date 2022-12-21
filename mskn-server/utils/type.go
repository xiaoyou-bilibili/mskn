package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

func Byte2Interface(data []byte) interface{} {
	var res interface{}
	err := json.Unmarshal(data, &res)
	if err != nil {
		log.Printf("unmarshal data err %v", err)
	}
	return res
}

func Interface2byte(data interface{}) []byte {
	tmp, err := json.Marshal(data)
	if err != nil {
		log.Printf("unmarshal data err %v", err)
	}
	return tmp
}

func Interface2String(data interface{}) string {
	return string(Interface2byte(data))
}

func String2Int64(data string) int64 {
	num, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
