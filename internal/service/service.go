package service

import (
	"encoding/json"
	"fmt"
	"kprg/internal/enricher"
	"kprg/internal/logger"
	"kprg/internal/models"
	"kprg/internal/producer"
	"kprg/internal/repository"
	"math/rand"
	"os"

	"github.com/brianvoe/gofakeit/v6"
)

func GetPostgresConnction() *models.Connection {

	connection := models.Connection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	return &connection

}

func GetKafkaConnction() *models.Connection {

	connection := models.Connection{
		Host: os.Getenv("KAFKA_HOST"),
		Port: os.Getenv("KAFKA_PORT"),
	}

	return &connection

}

func GetRedisConnction() *models.Connection {

	connection := models.Connection{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}

	return &connection

}

func SendRandomUser(logger logger.Logger, producer *producer.Producer, topic string) {

	rnum := rand.Intn(10)

	userCreateRequest := models.UserCreateRequest{}

	if rnum%7 != 0 {
		userCreateRequest.Name = gofakeit.FirstName()
	}

	if rnum%5 != 0 {
		userCreateRequest.Surname = gofakeit.LastName()
	}

	if rnum%3 != 0 {
		userCreateRequest.Patronymic = gofakeit.MiddleName()
	}

	message, err := json.Marshal(userCreateRequest)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal fio: %s", err))
		return
	}

	err = producer.Send(topic, message)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to produce fio: %s", err))
		return
	}

	logger.Info(fmt.Sprintf("Sent message: %s", string(message)))

}

func KafkaCreateUser(enricher *enricher.Enricher, repository repository.Repository, producer *producer.Producer, topic string, message []byte) error {

	userCreateRequest := models.UserCreateRequest{}

	err := json.Unmarshal(message, &userCreateRequest)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal message: %w", err)
	}

	errors := CheckUserCreateRequest(userCreateRequest)

	if len(errors) > 0 {

		userFailedResponse := models.UserFailedResponse{
			Name:       userCreateRequest.Name,
			Surname:    userCreateRequest.Surname,
			Patronymic: userCreateRequest.Patronymic,
			Errors:     errors,
		}

		message, err = json.Marshal(userFailedResponse)
		if err != nil {
			return fmt.Errorf("Failed to marshal FioFailed response: %w", err)
		}

		err = producer.Send(topic, message)
		if err != nil {
			return fmt.Errorf("Failed to produce FioFailed response: %w", err)
		}

		return nil

	}

	_, err = CreateUser(userCreateRequest, enricher, repository)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func CheckUserCreateRequest(request models.UserCreateRequest) []models.ErrorMessageResponse {

	errors := []models.ErrorMessageResponse{}

	if request.Name == "" {
		errors = append(errors, models.ErrorMessageResponse{Code: 1, Message: "Name is empty"})
	}

	if request.Surname == "" {
		errors = append(errors, models.ErrorMessageResponse{Code: 2, Message: "Surname is empty"})
	}

	return errors

}

func CreateUser(request models.UserCreateRequest, enricher *enricher.Enricher, repository repository.Repository) (*models.UserCreateResponse, error) {

	var err error

	user := models.User{
		Name:       request.Name,
		Surname:    request.Surname,
		Patronymic: request.Patronymic,
	}

	user.Age, err = enricher.GetAge(user.Name)
	if err != nil {
		return nil, fmt.Errorf("Failed to enrich user %s age: %s\n", user.Name, err)
	}

	user.Gender, err = enricher.GetGender(user.Name)
	if err != nil {
		return nil, fmt.Errorf("Failed to enrich user %s sex: %s\n", user.Name, err)
	}

	user.Nationality, err = enricher.GetNationality(user.Name)
	if err != nil {
		return nil, fmt.Errorf("Failed to enrich user %s nationality: %s\n", user.Name, err)
	}

	err = repository.CreateUser(&user)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user %+v: %s\n", user, err)
	}

	response := &models.UserCreateResponse{
		Id:          user.Id,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Gender:      user.Gender,
		Age:         user.Age,
		Nationality: user.Nationality,
	}

	return response, nil

}

func CheckUserUpdateRequest(request models.UserUpdateRequest) []models.ErrorMessageResponse {

	errors := []models.ErrorMessageResponse{}

	if request.Name == "" {
		errors = append(errors, models.ErrorMessageResponse{Code: 1, Message: "Name is empty"})
	}

	if request.Surname == "" {
		errors = append(errors, models.ErrorMessageResponse{Code: 2, Message: "Surname is empty"})
	}

	return errors

}

func UpdateUser(id string, request models.UserUpdateRequest, enricher *enricher.Enricher, repository repository.Repository) (*models.UserUpdateResponse, error) {

	user, err := repository.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user %s: %s\n", id, err)
	}

	if user == nil {
		return nil, fmt.Errorf("User Id %s not found", id)
	}

	if user.Name != request.Name {

		user.Name = request.Name

		user.Age, err = enricher.GetAge(user.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to enrich user %s age: %s\n", user.Name, err)
		}

		user.Gender, err = enricher.GetGender(user.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to enrich user %s sex: %s\n", user.Name, err)
		}

		user.Nationality, err = enricher.GetNationality(user.Name)
		if err != nil {
			return nil, fmt.Errorf("Failed to enrich user %s nationality: %s\n", user.Name, err)
		}

	}

	user.Surname = request.Surname
	user.Patronymic = request.Patronymic

	err = repository.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to update user %+v: %s\n", user, err)
	}

	response := &models.UserUpdateResponse{
		Id:          user.Id,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Gender:      user.Gender,
		Age:         user.Age,
		Nationality: user.Nationality,
	}

	return response, nil

}
