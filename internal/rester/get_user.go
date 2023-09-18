package rester

import (
	"fmt"
	"kprg/internal/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *Rester) GetUser(c *fiber.Ctx) error {

	id := c.Params("id", "")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON("User Id is not correct")
	}

	user, err := r.repository.GetUser(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("Failed to get user %s: %w", id, err)
	}

	if user == nil {
		return c.Status(http.StatusNotFound).JSON(fmt.Sprintf("User Id %s not found", id))
	}

	response :=
		models.UserGetResponse{
			Id:          user.Id,
			Name:        user.Name,
			Surname:     user.Surname,
			Patronymic:  user.Patronymic,
			Gender:      user.Gender,
			Age:         user.Age,
			Nationality: user.Nationality,
		}

	return c.Status(http.StatusOK).JSON(response)

}
