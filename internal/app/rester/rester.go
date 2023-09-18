package rester

import (
	"kprg/internal/cacher"
	"kprg/internal/enricher"
	"kprg/internal/repository"
	"kprg/internal/rester"
	"kprg/internal/service"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Start() {

	godotenv.Load(".env")

	rester_port := os.Getenv("RESTER_PORT")

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

	log.Printf("Created enricher %+v\n", enricher)

	// Rester
	rester := rester.NewRester(repository, enricher)

	err = rester.Listen(rester_port)
	if err != nil {
		log.Fatalf("Failed to start rester: %s\n", err)
	}

	defer rester.Close()

	log.Printf("Created rester %+v\n", rester)

}
