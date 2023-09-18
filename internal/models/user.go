package models

import "time"

type User struct {
	Id          string `gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name        string
	Surname     string
	Patronymic  string
	Gender      string
	Age         int
	Nationality string
	CreatedAt   time.Time
	ChangedAt   time.Time
	DeletedAt   time.Time
}

type UserCreateRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type UserCreateResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"sex"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type UserGetResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"sex"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type UsersGetResponse []UserGetResponse

type UserFailedResponse struct {
	Name       string                 `json:"name"`
	Surname    string                 `json:"surname"`
	Patronymic string                 `json:"patronymic"`
	Errors     []ErrorMessageResponse `json:"errors"`
}

type UserUpdateRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type UserUpdateResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"sex"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}
