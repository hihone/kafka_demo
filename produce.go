package kafka_demo

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func WriteMessage() error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicDemo,
		Balancer: &kafka.LeastBytes{}, // 负载均衡策略
		//Balancer: &kafka.Hash{}, // 基于 Key 的 Hash 分配到相同分区
	})

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		err := writer.WriteMessages(ctx, kafka.Message{
			HighWaterMark: 0,
			Key:           []byte("name"),
			Value:         []byte("Hihone"),
		}, kafka.Message{
			HighWaterMark: 0,
			Key:           []byte("name"),
			Value:         []byte("YYanghf"),
		})
		if err != nil {
			if errors.Is(err, kafka.UnknownTopicOrPartition) {
				continue
			}
			log.Fatal("发送消息出错啦, Error:", err)
		}
		fmt.Println("发送消息成功")
	}
	if err := writer.Close(); err != nil {
		log.Fatal("关闭发送消息失败", err)
	}

	return nil
}
