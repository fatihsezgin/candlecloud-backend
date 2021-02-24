package storage

type Store interface {
	Users() UserRepository
}
