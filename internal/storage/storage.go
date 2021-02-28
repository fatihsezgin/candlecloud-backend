package storage

type Store interface {
	Users() UserRepository
	Products() ProductRepository
}
