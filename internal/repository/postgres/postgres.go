package postgres

import (
	"fmt"
	"kprg/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	client *gorm.DB
}

func NewPostgres(connection *models.Connection) (*Postgres, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable TimeZone=Europe/Moscow", connection.Host, connection.Port, connection.User, connection.Password, connection.Database)

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("faild to connect to db: %w", err)
	}

	postgres :=
		&Postgres{
			client: client,
		}

	return postgres, nil

}

func (p *Postgres) AutoMigrate() error {

	err := p.autoMigrateUser()
	if err != nil {
		return err
	}

	return nil

}

func (p *Postgres) Close() error {

	db, err := p.client.DB()
	if err != nil {
		return fmt.Errorf("faild to close connection to db: %w", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("faild to close connection to db: %w", err)
	}

	return nil

}
