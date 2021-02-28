package storage

import (
	"github.com/fatihsezgin/candlecloud-backend/model"
)

// UserRepository interface is the common interface for a repository
// Each method checks the entity type.
type UserRepository interface {
	// All returns all the data in the repository.
	All() ([]model.User, error)
	// FindAll returns the entities matching the arguments.
	FindAll(argsStr map[string]string, argsInt map[string]int) ([]model.User, error)
	// FindByID finds the entity regarding to its ID.
	FindByID(id uint) (*model.User, error)
	// FindByUUID finds the entity regarding to its UUID.
	FindByUUID(uuid string) (*model.User, error)
	// FindByEmail finds the entity regarding to its Email.
	FindByEmail(email string) (*model.User, error)
	// FindByCredentials finds the entity regarding to its Email and Master Password.
	FindByCredentials(email, masterPassword string) (*model.User, error)
	// Save stores the entity to the repository
	Save(login *model.User) (*model.User, error)
	// Delete removes the entity from the store
	Delete(id uint, schema string) error
	// Migrate migrates the repository
	Migrate() error
	// CreateSchema creates schema for user
	CreateSchema(schema string) error
}

type ProductRepository interface {
	// All returns all the data in the repository.
	All() ([]model.Product, error)
	// FindAll returns the entities matching the arguments.
	//FindAll(argsStr map[string]string, argsInt map[string]int) ([]model.Product, error)
	// FindByID finds the entity regarding to its ID.
	FindByID(id uint) (*model.Product, error)
	// FindByUUID finds the entity regarding to its UUID.
	FindByUUID(uuid string) (*model.Product, error)
	// Save stores the entity to the repository
	Save(product *model.Product) (*model.Product, error)
	// Migrate migrates the repository
	Migrate() error
}
