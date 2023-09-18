package rester

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (r *Rester) DeleteUser(c *fiber.Ctx) error {

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

	err = r.repository.DeleteUser(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("Failed to delete user %s: %w", id, err)
	}

	c.Status(http.StatusOK)

	return nil

}
