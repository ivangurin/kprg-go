package grapher

import (
	"encoding/json"
	"fmt"
	"kprg/internal/enricher"
	"kprg/internal/repository"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Grapher struct {
	repository   repository.Repository
	enricher     *enricher.Enricher
	schema       graphql.Schema
	userType     *graphql.Object
	queryType    *graphql.Object
	mutationType *graphql.Object
}

func NewGrapher(repository repository.Repository, enricher *enricher.Enricher) (*Grapher, error) {

	var err error

	grapher := &Grapher{
		repository: repository,
		enricher:   enricher,
	}

	grapher.userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"surname": &graphql.Field{
					Type: graphql.String,
				},
				"patronymic": &graphql.Field{
					Type: graphql.String,
				},
				"gender": &graphql.Field{
					Type: graphql.String,
				},
				"age": &graphql.Field{
					Type: graphql.Int,
				},
				"nationality": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	grapher.queryType =
		graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{

					/* Get (read) single user by id
					   http://localhost:8080/users?query={user(id:"1"){id,name,surname, patronymic, age, gender, nationality}}
					*/
					"user": &graphql.Field{
						Type:        grapher.userType,
						Description: "Get user by id",
						Args: graphql.FieldConfigArgument{
							"id": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},
						Resolve: grapher.GetUser,
					},

					/* Get (read) user list
					   http://localhost:8080/users?query={list{id,name,surname, patronymic, age, gender, nationality}}
					*/
					"list": &graphql.Field{
						Type:        graphql.NewList(grapher.userType),
						Description: "Get user list",
						Resolve:     grapher.GetUsers,
					},
				},
			})

	grapher.mutationType =
		graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{

				/* Create new user
				http://localhost:8080/users?query=mutation+_{create(name:"Inca Kola"){id,name}}
				*/
				"create": &graphql.Field{
					Type:        grapher.userType,
					Description: "Create new user",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"surname": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"patronymic": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: grapher.CreateUser,
				},

				/* Update product by id
				   http://localhost:8080/users?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
				*/
				"update": &graphql.Field{
					Type:        grapher.userType,
					Description: "Update user by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"surname": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"patronymic": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: grapher.UpdateUser,
				},

				/* Delete product by id
				   http://localhost:8080/users?query=mutation+_{delete(id:1){id,name}}
				*/
				"delete": &graphql.Field{
					Type:        grapher.userType,
					Description: "Delete user by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: grapher.DeleteUser,
				},
			},
		})

	grapher.schema, err =
		graphql.NewSchema(
			graphql.SchemaConfig{
				Query:    grapher.queryType,
				Mutation: grapher.mutationType,
			},
		)

	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return grapher, nil

}

func (g *Grapher) Query(query string) *graphql.Result {

	result :=
		graphql.Do(graphql.Params{
			Schema:        g.schema,
			RequestString: query,
		})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}

	return result

}

func (g *Grapher) Listen(port string) error {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		result := g.Query(r.URL.Query().Get("query"))
		json.NewEncoder(w).Encode(result)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		return fmt.Errorf("failed to listen port %s: %w", port, err)
	}

	return nil

}

func (g *Grapher) Close() error {

	return nil

}
