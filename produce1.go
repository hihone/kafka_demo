package kafka_demo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func WriteMessage1() error {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", TopicDemo, Partition0)
	if err != nil {
		log.Fatal("写消息链接失败，Error：", errors.WithStack(err))
	}

	if err = conn.SetWriteDeadline(time.Now().UTC().Add(15 * time.Second)); err != nil {
		log.Fatal("写消息设置过期时间失败，Error：", errors.WithStack(err))
	}
	n, err := conn.WriteMessages(kafka.Message{
		Key:   []byte("name"),
		Value: []byte("Hihone"),
	}, kafka.Message{
		Key:   []byte("name"),
		Value: []byte("YYanghf"),
	})
	if err != nil {
		log.Fatal("发送消息写入失败，Error：", errors.WithStack(err))
	}
	fmt.Println("n ====", n)

	if err := conn.Close(); err != nil {
		log.Fatal("关闭发送消息失败", errors.WithStack(err))
	}

	fmt.Println("发送消息成功")

	return nil
}
