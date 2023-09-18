package consumer

import (
	"fmt"
	"kprg/internal/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	Concumer *kafka.Consumer
}

func NewConsumer(connection *models.Connection, groupID string) (*Consumer, error) {

	servers := fmt.Sprintf("%s:%s", connection.Host, connection.Port)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        servers,
		"broker.address.family":    "v4",
		"group.id":                 groupID,
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{Concumer: consumer}, nil

}

func (c *Consumer) Close() {

	c.Concumer.Close()

}
