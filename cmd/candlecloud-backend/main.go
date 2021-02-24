package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatihsezgin/candlecloud-backend/internal/app"
	"github.com/fatihsezgin/candlecloud-backend/internal/config"
	"github.com/fatihsezgin/candlecloud-backend/internal/router"
	"github.com/fatihsezgin/candlecloud-backend/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "[candle-cloud] ", 0)

	cfg, err := config.SetupConfigDefaults()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.DBConn(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New(db)
	app.MigrateSystemTables(s)

	srv := &http.Server{
		MaxHeaderBytes: 10, // 10 MB
		Addr:           ":" + cfg.Server.Port,
		WriteTimeout:   time.Second * time.Duration(cfg.Server.Timeout),
		ReadTimeout:    time.Second * time.Duration(cfg.Server.Timeout),
		IdleTimeout:    time.Second * 60,
		Handler:        router.New(s),
	}

	logger.Printf("listening on %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
