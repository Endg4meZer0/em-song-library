package cmd

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"effective-mobile-song-library/config"
	"effective-mobile-song-library/internal/delivery/http"
	"effective-mobile-song-library/internal/repository/api"
	pgDB "effective-mobile-song-library/internal/repository/db"
	"effective-mobile-song-library/internal/service"
	"effective-mobile-song-library/pkg/logger"
	_ "effective-mobile-song-library/docs"
)

// @title Song Library API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.PrintError(err, nil)
	}

	// Connect to DB
	db, err := openDB(*cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	// Database migrations
	migrationDriver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", migrationDriver)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.PrintFatal(err, nil)
	}

	// prepare repo
	carRepo := pgDB.NewSongsRepository(db)
	peopleRepo := pgDB.NewSongDetailsRepository(db)
	apiClient := api.NewApiClient(cfg)

	// service layer
	songLibraryService := service.NewSongLibraryService(carRepo, peopleRepo, apiClient)

	// handler
	handler := http.NewHandler(songLibraryService)

	srv := NewServer(
		handler,
		cfg,
	)

	err = srv.Start()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}