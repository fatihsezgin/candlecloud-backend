package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User Model
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UUID      uuid.UUID  `gorm:"type:uuid; type:varchar(100);"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Email     string     `json:"email" encrypt:"true"`
	Password  string     `json:"password" encrypt:"true"`
	Schema    string     `json:"schema"`
	Role      string     `json:"role"`
}

// UserDTO object for User type
type UserDTO struct {
	ID       uint      `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name" validate:"max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,max=100,min=6"`
	Schema   string    `json:"schema"`
	Role     string    `json:"role"`
}

// UserSignup object for Auth Signup endpoint
type UserSignup struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=100,min=6"`
}

type UserDTOTable struct {
	ID     uint      `json:"id"`
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Schema string    `json:"schema"`
	Role   string    `json:"role"`
}

func ConvertUserDTO(userSignup *UserSignup) *UserDTO {
	return &UserDTO{
		Email:    userSignup.Email,
		Password: userSignup.Password,
	}
}

func ToUser(userDTO *UserDTO) *User {
	return &User{
		ID:       userDTO.ID,
		UUID:     userDTO.UUID,
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: userDTO.Password,
		Schema:   userDTO.Schema,
		Role:     userDTO.Role,
	}
}

func ToUserDTO(user *User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		UUID:     user.UUID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Schema:   user.Schema,
		Role:     user.Role,
	}
}

func ToUserDTOTable(user User) UserDTOTable {
	return UserDTOTable{
		ID:     user.ID,
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Schema: user.Schema,
		Role:   user.Role,
	}
}

func ToUserDTOs(users []User) []UserDTOTable {
	userDTOs := make([]UserDTOTable, len(users))

	for i, item := range users {
		userDTOs[i] = ToUserDTOTable(item)
	}
	return userDTOs
}
