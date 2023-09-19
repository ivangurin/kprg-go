package grapher

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

func (g *Grapher) DeleteUser(params graphql.ResolveParams) (interface{}, error) {

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

	err = g.repository.DeleteUser(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete user %s exists: %w", id, err)
	}

	return nil, nil

}
