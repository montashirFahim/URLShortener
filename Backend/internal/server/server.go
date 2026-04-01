package server

import (
	"Server/internal/server/handler/api/v1"
	"Server/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	routes    Routes
	svc       *service.Service
	jwtSecret string
}

func NewServer(svc *service.Service, jwtSecret string) *Server {
	return &Server{
		svc:       svc,
		jwtSecret: jwtSecret,
	}
}

func (s *Server) addRoutes() {
	authHandler := v1.NewAuthHandler(s.svc.User)
	urlHandler := v1.NewURLHandler(s.svc.URL)

	// Public Auth Endpoints
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/api/v1/auth/register",
		Handler: http.HandlerFunc(authHandler.Register),
	})
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/oauth/token",
		Handler: http.HandlerFunc(authHandler.Token),
	})
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/oauth/revoke",
		Handler: http.HandlerFunc(authHandler.Revoke),
	})

	// Public Redirect Endpoint
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/{short_url}",
		Handler: http.HandlerFunc(urlHandler.Redirect),
	})

	// Private Endpoints (Wrapped with Auth Middleware)
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/api/v1/users/{id}/urls",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.Create)),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/users/{id}/urls",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.List)),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/users/{id}/urls/{url_id}",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.Get)),
	})
	s.routes.Add(Route{
		Method:  http.MethodDelete,
		Pattern: "/api/v1/users/{id}/urls/{url_id}",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.Delete)),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/users/{id}/urls/{url_id}/analytics",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.Analytics)),
	})
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/api/v1/gen",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.GenerateGuest)),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/gen/{short_url}",
		Handler: s.authMiddleware(http.HandlerFunc(urlHandler.GetGuestInfo)),
	})
}

func (s *Server) mux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	for _, route := range s.routes {
		r.Method(route.Method, route.Pattern, route.Handler)
	}

	return r
}

func (s *Server) Serve(port string) error {
	s.addRoutes()
	mux := s.mux()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	return server.ListenAndServe()
}
