package grapher

import (
	"kprg/internal/cacher"
	"kprg/internal/enricher"
	"kprg/internal/grapher"
	"kprg/internal/repository"
	"kprg/internal/service"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Start() {

	godotenv.Load(".env")

	grapher_port := os.Getenv("GRAPHER_PORT")

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

	// Grapher
	grapher := grapher.NewGrapher(repository, enricher)

	err = grapher.Listen(grapher_port)
	if err != nil {
		log.Fatalf("Failed to start grapher: %s\n", err)
	}

	defer grapher.Close()

	log.Printf("Created grapher %+v\n", grapher)

}
