package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Block struct {
	Height    int64    `json:"blockNumber"`
	BlockHash string   `json:"blockHash"`
	Txs       []string `json:"transactionList"`
}

func main() {
	// make a new reader that consumes from topic-A, partition 0, at last offset
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"127.0.0.1:9092"},
		Topic:     "block",
		Partition: 0,
		MinBytes:  256,  // 1KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(kafka.LastOffset)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			break
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		// 解析message，输出交易信息
		var block Block
		err = json.Unmarshal(m.Value, &block)
		if err != nil {
			fmt.Printf("error:%s\n", err.Error())
			break
		}
		fmt.Printf("hash:%s,height:%d,txs:%v\n", block.BlockHash, block.Height, block.Txs)
	}
}
