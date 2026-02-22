package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	httpserver "golang/internal/handler/http"
	"golang/internal/repository"
	"golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"
)

func Run() error {
	log.SetFlags(log.LstdFlags)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgCfg := initPostgreConfig()
	db := _postgres.NewPGDialect(ctx, pgCfg)
	defer func() { _ = db.DB.Close() }()

	repos := repository.NewRepositories(db)
	userUC := usecase.NewUserUsecase(repos.UserRepository)

	srvCfg := initServerConfig()
	router := httpserver.NewRouter(srvCfg, userUC)

	server := &http.Server{
		Addr:              srvCfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server started on %s", srvCfg.Addr)
	return server.ListenAndServe()
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        env("PG_HOST", "localhost"),
		Port:        env("PG_PORT", "5432"),
		Username:    env("PG_USER", "postgres"),
		Password:    env("PG_PASS", "postgres"),
		DBName:      env("PG_DB", "mydb"),
		SSLMode:     env("PG_SSL", "disable"),
		ExecTimeout: 5 * time.Second,
	}
}

func initServerConfig() *modules.ServerConfig {
	return &modules.ServerConfig{
		Addr:   env("APP_ADDR", ":8080"),
		APIKey: env("API_KEY", "dev-secret"),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
