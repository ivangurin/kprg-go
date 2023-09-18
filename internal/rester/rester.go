package rester

import (
	"fmt"
	"kprg/internal/enricher"
	"kprg/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Rester struct {
	repository repository.Repository
	enricher   *enricher.Enricher
	fiber      *fiber.App
}

func NewRester(repository repository.Repository, enricher *enricher.Enricher) *Rester {

	rester := &Rester{
		repository: repository,
		enricher:   enricher,
		fiber:      fiber.New(),
	}

	rester.fiber.Use(cors.New())

	rester.fiber.Get("/users", rester.GetUsers)
	rester.fiber.Get("/users/:id", rester.GetUser)
	rester.fiber.Post("/users", rester.CreateUser)
	rester.fiber.Post("/users/:id", rester.UpdateUser)
	rester.fiber.Delete("/users/:id", rester.DeleteUser)

	return rester

}

func (r *Rester) Listen(port string) error {

	err := r.fiber.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen port %s: %w", port, err)
	}

	return nil

}

func (r *Rester) Close() error {

	err := r.fiber.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown fiber: %w", err)
	}

	return nil

}
