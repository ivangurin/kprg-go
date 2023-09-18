package producer

import (
	"kprg/internal/producer"
	"kprg/internal/service"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Start() {

	godotenv.Load(".env")

	kafkaConnection := service.GetKafkaConnction()
	kafkaTopicFio := os.Getenv("KAFKA_TOPIC_FIO")

	producer, err := producer.NewProducer(kafkaConnection)
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}

	log.Printf("Created Producer %v\n", producer)

	defer producer.Close()

	ticker := time.NewTicker(time.Second)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run {
		select {
		case <-ticker.C:
			service.SendRandomUser(producer, kafkaTopicFio)

		case signal := <-shutdown:
			log.Printf("Caught signal %v: terminating\n", signal)
			run = false
		}
	}

}
