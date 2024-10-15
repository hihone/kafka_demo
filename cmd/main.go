package main

import "github.com/hihone/kafka_demo"

func main() {
	//_ = kafka_demo.WriteMessage1()
	kafka_demo.ReadMessage1()

	//read := kafka_demo.InitRead()
	//go kafka_demo.SignalLister(read)
	//kafka_demo.ReadMessage(read)
}
