package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hihone/kafka_demo/ibm_sarama_demo"
	"github.com/pkg/errors"
	"time"
)

func main() {
	readMessage()
}

func readMessage() {
	config := sarama.NewConfig()
	config.ClientID = "app"
	config.Consumer.Return.Errors = true
	config.Consumer.Fetch.Default = 3 * 1024 * 1024
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	c, err := sarama.NewConsumerGroup([]string{ibm_sarama_demo.ADDR}, ibm_sarama_demo.GROUP, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = c.Close()
	}()

	ctx := context.Background()
	for {
		err = c.Consume(ctx, []string{ibm_sarama_demo.TOPIC}, &ConsumerGroupHandle{})
		if err != nil {
			fmt.Printf("Error from consumer: %+v", errors.WithStack(err))
		}
	}
}

type ConsumerGroupHandle struct{}

func (h *ConsumerGroupHandle) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("消费者启动啦")
	return nil
}

func (h *ConsumerGroupHandle) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("消费者关闭啦")
	return nil
}

func (h *ConsumerGroupHandle) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println("Topic = ", msg.Topic, "，Partition = ", msg.Partition, "，Offset = ", msg.Offset, "，Key = ", string(msg.Key), "，Value = ", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}
