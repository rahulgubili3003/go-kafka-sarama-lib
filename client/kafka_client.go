package client

import "github.com/IBM/sarama"

func ConnectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	// Create new consumer
	conn, err := sarama.NewConsumer(brokersUrl, config)

	if err != nil {
		return nil, err
	}
	return conn, nil
}
