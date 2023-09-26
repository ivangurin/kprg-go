package producer

import (
	"fmt"
	"kprg/internal/logger"
	"kprg/internal/producer"
	"kprg/internal/service"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Start() {

	godotenv.Load(".env")

	logger := logger.NewLogger(os.Stdout)

	// Kafka producer
	kafkaConnection := service.GetKafkaConnction()
	kafkaTopicFio := os.Getenv("KAFKA_TOPIC_FIO")

	producer, err := producer.NewProducer(kafkaConnection)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create producer: %s", err))
		return
	}

	defer producer.Close()

	logger.Info(fmt.Sprintf("Created Producer %v", producer))

	ticker := time.NewTicker(time.Second)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run {
		select {
		case <-ticker.C:
			service.SendRandomUser(logger, producer, kafkaTopicFio)
			logger.Info(fmt.Sprintf("Sent random user to topic %v", kafkaTopicFio))

		case signal := <-shutdown:
			run = false
			logger.Info(fmt.Sprintf("Caught signal %v: terminating", signal))
		}
	}

}
