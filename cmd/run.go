package cmd

import (
	"flag"
	"log"
	"log/slog"
	"millionaire-list/config"
	"millionaire-list/internal/handler"
	"millionaire-list/internal/logger"
	"millionaire-list/internal/repo"
	"millionaire-list/internal/router"
	"millionaire-list/internal/service"
	"millionaire-list/migrations"
)

func Run() {
	migrationDirection := flag.String("migrate", "", "Run migration in direction: up or down")
	flag.Parse()

	cfg, err := config.InitConfig(".env")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	logger := logger.SetupLogger(cfg.Env) // предполагается, что cfg.Env содержит "local", "dev" или "prod"
	logger.Info("Config loaded", slog.Any("config", cfg))

	db, err := config.ConnectDB(cfg)
	if err != nil {
		logger.Error("Could not connect to database", slog.String("error", err.Error()))
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	logger.Info("Successfully connected to the database")

	if *migrationDirection != "" {
		logger.Info("Running migrations", slog.String("direction", *migrationDirection))
		switch *migrationDirection {
		case "up":
			logger.Info("Running migration UP...")
			if err := migrations.RunMigrationUp(db); err != nil {
				logger.Error("Migration up failed", slog.String("error", err.Error()))
				log.Fatalf("Migration up failed: %v", err)
			}
			logger.Info("Migration up finished")
		case "down":
			logger.Info("Running migration DOWN...")
			if err := migrations.RunMigrationDown(db); err != nil {
				logger.Error("Migration down failed", slog.String("error", err.Error()))
				log.Fatalf("Migration down failed: %v", err)
			}
			logger.Info("Migration down finished")
		default:
			logger.Error("Invalid migration direction", slog.String("direction", *migrationDirection))
			log.Fatalf("Invalid migration direction: %s. Use 'up' or 'down'.", *migrationDirection)
		}
		return
	}

	millionaireRepo := repo.NewMillionaireRepo(db)
	millionaireService := service.NewMillionaireService(millionaireRepo)
	millionaireHandler := handler.NewMillionaireHandler(millionaireService)

	r := router.SetupRouter(millionaireHandler)

	// Вывод всех зарегистрированных маршрутов
	log.Println("Registered routes:")
	for _, route := range r.Routes() {
		logger.Info("Route registered", slog.Any("route", route))
	}

	// Запуск сервера
	if err := r.Run(); err != nil {
		logger.Error("Failed to start server", slog.String("error", err.Error()))
		log.Fatalf("Failed to start server: %v", err)
	}
}
