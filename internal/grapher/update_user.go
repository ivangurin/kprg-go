package grapher

import (
	"errors"
	"fmt"
	"kprg/internal/models"
	"kprg/internal/service"

	"github.com/graphql-go/graphql"
)

func (g *Grapher) UpdateUser(params graphql.ResolveParams) (interface{}, error) {

	id, _ := params.Args["id"].(string)
	if id == "" {
		return nil, errors.New("User Id is not correct")
	}

	exists, err := g.repository.IsUserExists(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to check user %s exists: %w", id, err)
	}

	if !exists {
		return nil, fmt.Errorf("User Id %s not found", id)
	}

	request := models.UserUpdateRequest{
		Name:       params.Args["name"].(string),
		Surname:    params.Args["surname"].(string),
		Patronymic: params.Args["patronymic"].(string),
	}

	errs := service.CheckUserUpdateRequest(request)

	if len(errs) > 0 {
		return nil, errors.New(errs[0].Message)
	}

	response, err := service.UpdateUser(id, request, g.enricher, g.repository)
	if err != nil {
		return nil, fmt.Errorf("failed to update user %s: %w", id, err)
	}

	return response, nil

}
