package v1

import (
	apimv "denet-test-task/internal/api/v1/middlewares"
	service "denet-test-task/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(r chi.Router, services *service.Services) {
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(apimv.SlogRequestContext)
	r.Use(apimv.SlogAccessLogger)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	// r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(cr chi.Router) {
		newAuthRoutes(cr, services.Auth)
	})

	authMiddleware := &apimv.AuthMiddleware{AuthService: services.Auth}
	r.Route("/api/v1", func(api chi.Router) {
		api.Use(authMiddleware.UserIdentity)
		// Остальные группы хендлеров
	})
}
