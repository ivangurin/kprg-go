package consumer

import (
	"kprg/internal/cacher"
	"kprg/internal/consumer"
	"kprg/internal/enricher"
	"kprg/internal/producer"
	"kprg/internal/repository"
	"kprg/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func Start() {

	godotenv.Load(".env")

	// DB
	postgresConnection := service.GetPostgresConnction()

	repository, err := repository.NewRepository(repository.DatabasePostgres, postgresConnection)
	if err != nil {
		log.Fatalf("Failed to create db connection: %s\n", err)
	}

	defer repository.Close()

	log.Printf("Created Repository %+v\n", repository)

	// Redis
	redisConnection := service.GetRedisConnction()

	cacher, err := cacher.NewCacher(redisConnection)
	if err != nil {
		log.Fatalf("Failed to create redis connection: %s\n", err)
	}

	defer cacher.Close()

	log.Printf("Created Cacher %+v\n", cacher)

	// Enricher
	enricher := enricher.NewEnricher(cacher)

	// Kafka
	kafkaTopicFio := os.Getenv("KAFKA_TOPIC_FIO")
	kafkaTopicFioFailed := os.Getenv("KAFKA_TOPIC_FIO_FAILED")

	kafkaConnection := service.GetKafkaConnction()

	// Producer
	producer, err := producer.NewProducer(kafkaConnection)
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	defer producer.Close()

	log.Printf("Created Producer %+v\n", producer)

	// Consumer
	consumer, err := consumer.NewConsumer(kafkaConnection, "Group1")
	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	defer consumer.Close()

	log.Printf("Created Consumer %+v\n", consumer)

	// Subsribe to topic
	err = consumer.Concumer.SubscribeTopics([]string{kafkaTopicFio}, nil)
	if err != nil {
		log.Fatalf("Failed to subsribe topic: %s\n", err)
	}

	// Listen system signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run {

		select {
		case signal := <-shutdown:
			log.Printf("Caught signal %v: terminating...\n", signal)
			run = false

		default:

			event := consumer.Concumer.Poll(100)
			if event == nil {
				continue
			}

			switch e := event.(type) {
			case *kafka.Message:

				log.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))

				err = service.KafkaCreateUser(enricher, repository, producer, kafkaTopicFioFailed, e.Value)
				if err != nil {
					log.Printf("Failed to save message %s:\n", err)
				}

				_, err = consumer.Concumer.StoreMessage(e)
				if err != nil {
					log.Printf("Failed to store offset after message %s:\n", e.TopicPartition)
				}

			case kafka.Error:
				log.Printf("Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}

			default:
				log.Printf("Ignored %v\n", e)
			}

		}

	}

	log.Printf("Shutdowing...\n")

}
