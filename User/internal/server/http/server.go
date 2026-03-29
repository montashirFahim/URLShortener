package http

import (
	"User/internal/server/http/handler/api/v1"
	"User/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	routes    Routes
	db        *sqlx.DB
	cache     *redis.Client
	svc       *service.Service
	jwtSecret []byte
}

func (s *Server) addRoutes() {
	userHandler := v1.NewUserHandler(s.svc.User, s.svc.Url, s.jwtSecret)

	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/api/v1/user/register",
		Handler: http.HandlerFunc(userHandler.Register),
	})
	s.routes.Add(Route{
		Method:  http.MethodPost,
		Pattern: "/api/v1/user/login",
		Handler: http.HandlerFunc(userHandler.Login),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/user/login",
		Handler: http.HandlerFunc(userHandler.LoginCheck),
	})
	s.routes.Add(Route{
		Method:  http.MethodGet,
		Pattern: "/api/v1/user/{userid}",
		Handler: http.HandlerFunc(userHandler.GetUserUrls),
	})
}

// NewServer returns instance of Server
func NewServer(db *sqlx.DB, cache *redis.Client, svc *service.Service, jwtSecret []byte) *Server {
	return &Server{
		db:        db,
		cache:     cache,
		svc:       svc,
		jwtSecret: jwtSecret,
	}
}

func (s *Server) mux() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(DetailedLogger) // Commenting out if not defined or failing

	for _, route := range s.routes {
		r.Method(route.Method, route.Pattern, route.Handler)
	}

	return r
}

func (s *Server) Serve() error {
	// monitoring := newAPM("ratings").Histogram(nil).Do
	s.addRoutes()
	// s.routes.ApplyRouteModifire(
	// 	//addNewrelic(s.NewRelic),
	// 	addSentry(&s.CacheContainer.Cnf.Sentry),
	// 	addHubble(monitoring),
	// )
	mux := s.mux()
	// mux.Mount("/metrics", promhttp.Handler())

	// Create HTTP server with timeout configurations
	server := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,  // Maximum duration for reading the entire request
		WriteTimeout: 60 * time.Second,  // Maximum duration before timing out writes of the response
		IdleTimeout:  120 * time.Second, // Maximum amount of time to wait for the next request when keep-alives are enabled
	}

	log.Println("start to listen on port: 8000")
	return server.ListenAndServe()
}
