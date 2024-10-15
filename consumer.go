package kafka_demo

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func InitRead() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicDemo,
		GroupID:  "demo-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func ReadMessage(read *kafka.Reader) {
	for {
		msg, err := read.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("读取信息出错啦，Error:", err)
			break
		}

		fmt.Printf("Topic: %s, Offset: %d, Key: %s, Msg: %s", msg.Topic, msg.Offset, msg.Key, msg.Value)
	}

}

func SignalLister(read *kafka.Reader) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	if read != nil {
		if err := read.Close(); err != nil {
			log.Fatal("关闭消费信息出错，Error：", err)
		}
		os.Exit(0)
	}
}
