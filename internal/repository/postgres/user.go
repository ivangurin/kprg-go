package postgres

import (
	"fmt"
	"kprg/internal/models"
	"time"
)

func (p *Postgres) autoMigrateUser() error {

	err := p.client.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to automigrate user: %s", err.Error())
	}

	return nil

}

func (p *Postgres) IsUserExists(id string) (bool, error) {

	user := &models.User{}

	result := p.client.Where("id = ?", id).Limit(1).Find(user)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check user exist id = %v: %s", id, result.Error.Error())
	}

	return result.RowsAffected > 0, nil

}

func (p *Postgres) CreateUser(user *models.User) error {

	user.CreatedAt = time.Now()

	err := p.client.Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user %+v: %s", user, err.Error())
	}

	return nil

}

func (p *Postgres) GetUser(id string) (*models.User, error) {

	user := models.User{}

	err := p.client.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to select user id = %s: %s", id, err.Error())
	}

	return &user, nil

}

func (p *Postgres) UpdateUser(user *models.User) error {

	user.ChangedAt = time.Now()

	err := p.client.Updates(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user %+v: %s", user, err.Error())
	}

	return nil

}

func (p *Postgres) DeleteUser(id string) error {

	user, err := p.GetUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %s", id, err.Error())
	}

	if user == nil {
		return fmt.Errorf("user %s does not exist: %s", id, err.Error())
	}

	user.DeletedAt = time.Now()

	err = p.client.Updates(user).Error
	if err != nil {
		return fmt.Errorf("failed to delete user %+v: %s", user, err.Error())
	}

	return nil

}

func (p *Postgres) GetUsers(where string, args []interface{}, offset int, limit int) ([]*models.User, error) {

	users := []*models.User{}

	query := p.client.Model(&models.User{})

	if where != "" {
		where += " AND "
	}

	where += "deleted_at = '0001-01-01 00:00:00+00'"

	query.Where(where, args...)

	if offset > 0 {
		query.Offset(int(offset))
	}

	if limit > 0 {
		query.Limit(int(limit))
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil

}
