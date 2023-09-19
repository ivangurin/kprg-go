package grapher

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func (r *Grapher) GetUser(p graphql.ResolveParams) (interface{}, error) {

	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, nil
	}

	user, err := r.repository.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id %s: %w", id, err)
	}

	if user != nil {
		return user, nil
	}

	return nil, nil

}
