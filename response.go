package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

func (r Response) ToString() string {
	// marshal the response into a JSON string
	jsonString, err := json.Marshal(r)
	if err != nil {
		return `{"success":false,"error":"Failed to marshal response","data":{}}`
	}
	return string(jsonString)
}

func PrintError(err string) {
	response := Response{Success: false, Error: err, Data: nil}
	fmt.Println(response.ToString())
	os.Exit(1)
}

func PrintData(data interface{}) {
	response := Response{Success: true, Error: "", Data: data}
	fmt.Println(response.ToString())
	os.Exit(0)
}
