package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	read  *kafka.Reader
	topic = "demo1"
)

func writeMessage(ctx context.Context) {
	writer := kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
		Transport:              nil,
		AllowAutoTopicCreation: true,
	}
	defer func() {
		_ = writer.Close()
	}()

	for i := 0; i < 3; i++ {
		err := writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte("1"),
			Value: []byte("111"),
		}, kafka.Message{
			Key:   []byte("2"),
			Value: []byte("222"),
		}, kafka.Message{
			Key:   []byte("3"),
			Value: []byte("333"),
		})
		if err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				log.Fatal("发送信息出错啦，Error：", errors.WithStack(err))
			}
		} else {
			break
		}
	}
}

func readMessage(ctx context.Context) {
	read = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupID:        "demo1_group",
		Topic:          topic,
		Partition:      0,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.LastOffset,
	})

	for {
		if msg, err := read.ReadMessage(ctx); err != nil {
			log.Fatal("读取消息出错啦，Error：", errors.WithStack(err))
		} else {
			fmt.Printf("Topic: %s, Partition: %d, Offset: %d, Key: %s, Value: %s \n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		}
	}
}

func listenSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	fmt.Println("接收到信号", sig.String())
	if read != nil {
		_ = read.Close()
		//if err := read.Close(); err != nil {
		//	log.Fatal("读消息关闭出错啦，Error：", errors.WithStack(err))
		//}
	}
	os.Exit(0)
}

func main() {
	ctx := context.Background()
	//writeMessage(ctx)

	go listenSignal()
	readMessage(ctx)
}
