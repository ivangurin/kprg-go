package rester

import (
	"fmt"
	"kprg/internal/models"
	"kprg/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *Rester) UpdateUser(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON("User Id is not correct")
	}

	exists, err := r.repository.IsUserExists(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("Failed to check user %s exists: %w", id, err)
	}

	if !exists {
		return c.Status(http.StatusNotFound).JSON(fmt.Sprintf("User Id %s not found", id))
	}

	request := models.UserUpdateRequest{}

	err = c.BodyParser(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return fmt.Errorf("failed to parse body: %w", err)
	}

	errors := service.CheckUserUpdateRequest(request)

	if len(errors) > 0 {

		userFailedResponse := models.UserFailedResponse{
			Name:       request.Name,
			Surname:    request.Surname,
			Patronymic: request.Patronymic,
			Errors:     errors,
		}

		return c.Status(http.StatusBadRequest).JSON(userFailedResponse)

	}

	response, err := service.UpdateUser(id, request, r.enricher, r.repository)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return c.Status(http.StatusOK).JSON(response)

}
