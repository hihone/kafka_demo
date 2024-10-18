package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hihone/kafka_demo/ibm_sarama_demo"
	"github.com/pkg/errors"
	"log"
	"time"
)

func main() {
	sendMessage()
	asyncSendMessage()
}

func sendMessage() {
	config := sarama.NewConfig()
	config.ClientID = "app"
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	p, err := sarama.NewSyncProducer([]string{ibm_sarama_demo.ADDR}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = p.Close()
	}()

	// 发送一条信息
	if _, _, err = p.SendMessage(&sarama.ProducerMessage{
		Topic: ibm_sarama_demo.TOPIC,
		Key:   sarama.StringEncoder("OnlyOne"),
		Value: sarama.StringEncoder("仅仅是发送一条消息"),
	}); err != nil {
		log.Fatal("发送一条消息出错啦，Error", errors.WithStack(err))
	}

	// 批量发送
	messages := make([]*sarama.ProducerMessage, 0)
	for i, s := range []string{"Hello", "Hi", "你好"} {
		messages = append(messages, &sarama.ProducerMessage{
			Topic: ibm_sarama_demo.TOPIC,
			Key:   sarama.StringEncoder(fmt.Sprintf("Key_%d", i)),
			Value: sarama.StringEncoder(s),
		})
	}
	if err = p.SendMessages(messages); err != nil {
		log.Fatal("发送信息出错啦, Error:", errors.WithStack(err))
	}
	fmt.Println("发送信息完成")
}

func asyncSendMessage() {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	p, err := sarama.NewAsyncProducer([]string{ibm_sarama_demo.ADDR}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = p.Close()
	}()

	go func() {
		for {
			select {
			case <-p.Successes():
				fmt.Println("")
				fmt.Println("异步发送消息成功")
			case err = <-p.Errors():
				log.Println("异步发送消息出错啦，Error：", errors.WithStack(err))

			}
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: ibm_sarama_demo.TOPIC,
		Key:   sarama.StringEncoder("OnlyOne"),
		Value: sarama.StringEncoder("异步消息发送"),
	}

	p.Input() <- msg
	time.Sleep(1 * time.Second)
}
