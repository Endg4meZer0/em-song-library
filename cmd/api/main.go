package cmd

import (
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"effective-mobile-song-library/config"
	_ "effective-mobile-song-library/docs"
	"effective-mobile-song-library/internal/delivery/http"
	pgDB "effective-mobile-song-library/internal/repository/db"
	"effective-mobile-song-library/internal/repository/external"
	"effective-mobile-song-library/internal/service"
	"effective-mobile-song-library/pkg/logger"
)

// @title Song Library API
// @version 1.0
// @host localhost:8080
// @BasePath /
func Start() {
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
	songsRepo := pgDB.NewSongsRepository(db)
	apiClient := external.NewApiClient(cfg)

	// service layer
	songLibraryService := service.NewSongLibraryService(songsRepo, apiClient)

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
