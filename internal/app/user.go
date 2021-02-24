package app

import (
	"fmt"

	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
	"github.com/fatihsezgin/candlecloud-backend/model"
	uuid "github.com/satori/go.uuid"
)

var (
	//ErrGenerateSchema represents message for generating schema
	ErrGenerateSchema = fmt.Errorf("an error occured while genarating schema")
	//ErrCreateSchema represents message for creating schema
	ErrCreateSchema = fmt.Errorf("an error occured while creating the schema and tables")
)

func CreateUser(s storage.Store, userDTO *model.UserDTO) (*model.User, error) {
	var err error

	err = PayloadValidator(userDTO)
	if err != nil {
		return nil, err
	}
	// Hashing the password with Bcrypt
	userDTO.Password = NewBcrypt([]byte(userDTO.Password))

	// Generate new UUID for user
	userDTO.UUID = uuid.NewV4()

	// New user's role is Member (not Admin)
	userDTO.Role = "Member"

	createdUser, err := s.Users().Save(model.ToUser(userDTO))
	if err != nil {
		return nil, err
	}

	updatedUser, err := GenerateSchema(s, createdUser)
	if err != nil {
		return nil, ErrCreateSchema
	}
	// Create user schema and tables
	err = s.Users().CreateSchema(updatedUser.Schema)
	if err != nil {
		return nil, ErrCreateSchema
	}

	// Create user tables in user schema
	// TODO design the order and migrate the user later
	//MigrateUserTables(s, updatedUser.Schema)

	return createdUser, nil
}

// UpdateUser updates the user with the dto and applies the changes in the store
func UpdateUser(s storage.Store, user *model.User, userDTO *model.UserDTO, isAuthorized bool) (*model.User, error) {

	// TODO: Refactor the contents of updated user with a logical way
	if userDTO.Password != "" && NewBcrypt([]byte(userDTO.Password)) != user.Password {
		userDTO.Password = NewBcrypt([]byte(userDTO.Password))
	} else {
		userDTO.Password = user.Password
	}

	user.Name = userDTO.Name
	user.Email = userDTO.Email
	user.Password = userDTO.Password

	// This never changes
	user.Schema = fmt.Sprintf("user%d", user.ID)

	// Only Admin's can change role
	if isAuthorized {
		user.Role = userDTO.Role
	}

	updatedUser, err := s.Users().Save(user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

// GenerateSchema creates user schema and tables
func GenerateSchema(s storage.Store, user *model.User) (*model.User, error) {
	user.Schema = fmt.Sprintf("user%d", user.ID)
	savedUser, err := s.Users().Save(user)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}
