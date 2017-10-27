package main

import (
	"encoding/json"
)

type Message struct {
	Msg string  `json:"msg"`
}

func main() {
	jsonStr := "{\"msg\":\"hello\"}"

	jsonByte := []byte(jsonStr)
	var message Message = Message{}
	var msgInterface interface{}
	json.Unmarshal(jsonByte, &message)
	json.Unmarshal(jsonByte, &msgInterface)

	println(message.Msg)
	println(msgInterface)
}
