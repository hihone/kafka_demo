package kafka_demo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func ReadMessage1() {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", TopicDemo, Partition0)
	if err != nil {
		log.Fatal("读消息链接失败，Error：", errors.WithStack(err))
	}

	if err = conn.SetReadDeadline(time.Now().UTC().Add(10 * time.Second)); err != nil {
		log.Fatal("设置读消息生命周期失败，Error：", errors.WithStack(err))
	}

	batch := conn.ReadBatch(10e3, 1e6)
	b := make([]byte, 10e3)
	for {
		m, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:m]))
	}
	if err = batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err = conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
