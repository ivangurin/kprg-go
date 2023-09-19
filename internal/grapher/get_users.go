package grapher

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func (g *Grapher) GetUsers(params graphql.ResolveParams) (interface{}, error) {

	users, err := g.repository.GetUsers("", nil, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil

}
