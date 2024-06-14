package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	fmt.Println("Hello")
	app := fiber.New()
	api := app.Group("/api/v1") // /api

	api.Post("/comments", createComment)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}

func createComment(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(Comment)

	//  Parse body into comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Println(err)
		err := c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		if err != nil {
			return err
		}
		return err
	}

	brokersUrl := []string{"localhost:9092"}
	topic := "messages_topic"
	message, err := json.Marshal(cmt)

	producer, err := connectProducer(brokersUrl)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	err = c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
	if err != nil {
		err := c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func connectProducer(brokersUrl []string) (sarama.SyncProducer, error) {

	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
