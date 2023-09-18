package producer

import (
	"fmt"
	"kprg/internal/models"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Producer *kafka.Producer
}

func NewProducer(connection *models.Connection) (*Producer, error) {

	servers := fmt.Sprintf("%s:%s", connection.Host, connection.Port)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": servers})
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				log.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	return &Producer{producer}, nil

}

func (p *Producer) Send(topic string, message []byte) error {

	for {

		err := p.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		}, nil)

		if err != nil {

			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				time.Sleep(time.Second)
				continue
			}

			return err

		}

		break

	}

	return nil

}

func (p *Producer) Close() {

	for p.Producer.Flush(5) > 0 {
		log.Println("Still waiting to flush outstanding messages")
	}

	p.Producer.Close()

}
