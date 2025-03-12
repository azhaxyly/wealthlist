package cmd

import (
	"flag"
	"log/slog"
	"wealthlist/config"
	"wealthlist/internal/handler"
	"wealthlist/internal/logger"
	"wealthlist/internal/repo"
	"wealthlist/internal/router"
	"wealthlist/internal/service"
	"wealthlist/migrations"
)

func Run() {
	migrationDirection := flag.String("migrate", "", "Run migration in direction: up or down")
	flag.Parse()

	cfg, err := config.InitConfig(".env")
	if err != nil {
		slog.Error("Could not load config", slog.String("error", err.Error()))
		return
	}

	log := logger.SetupLogger(cfg.Env)
	log.Info("Config loaded", slog.Any("config", cfg))

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Error("Could not connect to database", logger.Err(err))
		return
	}
	defer db.Close()

	log.Info("Successfully connected to the database")

	if *migrationDirection != "" {
		log.Info("Running migrations", slog.String("direction", *migrationDirection))
		switch *migrationDirection {
		case "up":
			log.Info("Running migration UP...")
			if err := migrations.RunMigrationUp(db); err != nil {
				log.Error("Migration up failed", logger.Err(err))
				return
			}
			log.Info("Migration up finished")
		case "down":
			log.Info("Running migration DOWN...")
			if err := migrations.RunMigrationDown(db); err != nil {
				log.Error("Migration down failed", logger.Err(err))
				return
			}
			log.Info("Migration down finished")
		default:
			log.Error("Invalid migration direction", logger.Err(err))
			return
		}
		return
	} else {
		log.Info("No migration flag provided, running UP migrations by default...")
		if err := migrations.RunMigrationUp(db); err != nil {
			log.Error("Migration up failed", logger.Err(err))
			return
		}
	}

	millionaireRepo := repo.NewMillionaireRepo(db, log)
	photoRepo := repo.NewPhotoRepo(db, log)

	millionaireService := service.NewMillionaireService(millionaireRepo, log)
	homeService := service.NewHomeService(millionaireRepo, log)
	photoService := service.NewPhotoService(photoRepo, log)
	feedbackService := service.NewFeedbackService(cfg, log)

	millionaireHandler := handler.NewMillionaireHandler(millionaireService, log)
	homeHandler := handler.NewHomeHandler(homeService, log)
	photoHandler := handler.NewPhotoHandler(photoService, log)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService, log)

	r := router.SetupRouter(millionaireHandler, photoHandler, homeHandler, feedbackHandler)

	log.Info("Starting server on :8080")
	if err := r.Run(); err != nil {
		log.Error("Failed to start server", logger.Err(err))
	}
}
