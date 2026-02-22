package http

import (
	"net/http"

	"golang/internal/handler/http/handlers"
	appmw "golang/internal/handler/http/middleware"
	"golang/internal/usecase"
	"golang/pkg/modules"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *modules.ServerConfig, userUC usecase.UserService) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(appmw.RequestLogger)

	r.Get("/healthz", handlers.Health)

	r.Route("/users", func(rr chi.Router) {
		rr.Use(appmw.Auth(cfg.APIKey))

		h := handlers.NewUsersHandler(userUC)

		rr.Get("/", h.GetUsers)
		rr.Post("/", h.CreateUser)
		rr.Get("/{id}", h.GetUserByID)
		rr.Patch("/{id}", h.UpdateUser)
		rr.Delete("/{id}", h.DeleteUser)
	})

	return r
}
