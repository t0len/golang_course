package app

import (
	"log"
	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// Postgres
	pg, err := postgres.New(
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.DBName,
		cfg.PG.SSLMode,
	)
	if err != nil {
		log.Fatalf("postgres connect error: %v", err)
	}

	// Автомиграция
	if err := pg.Conn.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("automigrate error: %v", err)
	}

	// Layers
	l := logger.New()
	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// HTTP
	handler := gin.Default()
	v1.NewRouter(handler, userUseCase, l)

	log.Printf("Server started on :%s", cfg.HTTPPort)
	if err := handler.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
