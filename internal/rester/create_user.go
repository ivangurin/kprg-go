package rester

import (
	"fmt"
	"kprg/internal/models"
	"kprg/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *Rester) CreateUser(c *fiber.Ctx) error {

	request := models.UserCreateRequest{}

	err := c.BodyParser(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf("failed to parse body: %w", err)
	}

	errors := service.CheckUserCreateRequest(request)

	if len(errors) > 0 {

		userFailedResponse := models.UserFailedResponse{
			Name:       request.Name,
			Surname:    request.Surname,
			Patronymic: request.Patronymic,
			Errors:     errors,
		}

		return c.Status(http.StatusBadRequest).JSON(userFailedResponse)

	}

	response, err := service.CreateUser(request, r.enricher, r.repository)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return c.Status(http.StatusOK).JSON(response)

}
