package grapher

import (
	"errors"
	"fmt"
	"kprg/internal/models"
	"kprg/internal/service"

	"github.com/graphql-go/graphql"
)

func (g *Grapher) CreateUser(params graphql.ResolveParams) (interface{}, error) {

	request :=
		models.UserCreateRequest{
			Name:       params.Args["name"].(string),
			Surname:    params.Args["surname"].(string),
			Patronymic: params.Args["patronymic"].(string),
		}

	errs := service.CheckUserCreateRequest(request)

	if len(errs) > 0 {
		return nil, errors.New(errs[0].Message)
	}

	response, err := service.CreateUser(request, g.enricher, g.repository)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return response, nil

}
