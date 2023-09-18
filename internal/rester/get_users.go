package rester

import (
	"fmt"
	"kprg/internal/models"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	rqp "github.com/timsolov/rest-query-parser"
)

func (r *Rester) GetUsers(c *fiber.Ctx) error {

	url, _ := url.Parse(string(c.Request().URI().FullURI()))

	q, err := rqp.NewParse(url.Query(), rqp.Validations{
		"name:":        nil,
		"surname:":     nil,
		"patronymic:":  nil,
		"gender:":      nil,
		"age:int":      nil,
		"nationality:": nil,
	})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			models.ErrorMessageResponse{
				Code:    0,
				Message: err.Error(),
			})
	}

	users, err := r.repository.GetUsers(q.Where(), q.Args(), q.Offset, q.Limit)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return fmt.Errorf("failed to get users: %w", err)
	}

	response := make(models.UsersGetResponse, 0, len(users))

	for _, user := range users {

		userResponse :=
			models.UserGetResponse{
				Id:          user.Id,
				Name:        user.Name,
				Surname:     user.Surname,
				Patronymic:  user.Patronymic,
				Gender:      user.Gender,
				Age:         user.Age,
				Nationality: user.Nationality,
			}

		response = append(response, userResponse)

	}

	return c.Status(http.StatusOK).JSON(response)

}
