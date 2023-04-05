package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	lookBack  = 100.0 // days
	batchSize = 20.0  // days
)

type Message struct {
	Data []struct {
		Time  int64  `json:"time"`
		Price string `json:"priceUsd"`
	} `json:"data"`
}

func GetData(assetId string) []Message {
	times := int(math.Ceil(lookBack / batchSize))

	var data []Message
	for i := 1; i <= times; i++ {
		start := time.Now().AddDate(0, 0, -(i+1)*int(batchSize)).UnixMilli()
		end := time.Now().AddDate(0, 0, -i*int(batchSize)).UnixMilli()

		m := GetDataInRange(assetId, start, end)
		if m.Data == nil {
			continue
		}

		data = append(data, m)
	}

	return data
}

func GetDataInRange(assetId string, start int64, end int64) Message {
	resp, err := http.Get(fmt.Sprintf(
		"https://api.coincap.io/v2/assets/%s/history?interval=h1&start=%d&end=%d",
		assetId, start, end))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return Message{}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return Message{}
	}
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return Message{}
	}

	var m Message
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Printf("Error while unmarshalling: %v\n", err)
		return Message{}
	}
	return m
}
