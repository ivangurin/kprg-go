package repository

import (
	"fmt"
	"kprg/internal/models"
	"kprg/internal/repository/postgres"
)

type Repository interface {
	AutoMigrate() error
	Close() error
	User
}

type User interface {
	IsUserExists(id string) (bool, error)
	CreateUser(user *models.User) error
	GetUser(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
	GetUsers(where string, args []interface{}, offset int, limit int) ([]*models.User, error)
}

const DatabasePostgres = "postgres"

func NewRepository(database string, connection *models.Connection) (Repository, error) {

	var err error
	var repository Repository

	switch connection.Database {
	case DatabasePostgres:
		repository, err = postgres.NewPostgres(connection)
	default:
		err = fmt.Errorf("unknown database: %s", connection.Database)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	err = repository.AutoMigrate()
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return repository, nil

}
