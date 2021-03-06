package app

import (
	"log"

	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
)

// MigrateSystemTables runs auto migration for the system models (User),
// will only add missing fields won't delete/change current data in the store.
func MigrateSystemTables(s storage.Store) {
	if err := s.Users().Migrate(); err != nil {
		log.Println(err)
	}
	if err := s.Products().Migrate(); err != nil {
		log.Println(err)
	}
}

// func MigrateUserTables(s storage.Store, schema string) {}

// MigrateUserTables runs auto migration for user models in user schema,
// will only add missing fields won't delete/change current data in the store.
// func MigrateUserTables(s storage.Store, schema string) {
// 	if err := s.Logins().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// 	if err := s.CreditCards().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// 	if err := s.BankAccounts().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// 	if err := s.Notes().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// 	if err := s.Emails().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// 	if err := s.Servers().Migrate(schema); err != nil {
// 		log.Println(err)
// 	}
// }
