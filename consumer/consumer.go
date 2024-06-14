package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/rahulgubili3003/go-kafka-sarama/client"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Could not load properties")
	}

	topic := os.Getenv("TOPIC")

	brokersUrl := []string{"localhost:9092"}

	worker, err := client.ConnectConsumer(brokersUrl)
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, msg.Topic, string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}
}
