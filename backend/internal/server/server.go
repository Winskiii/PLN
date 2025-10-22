package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	appcfg "backend/internal/config"
	apphandlers "backend/internal/handlers"
	appmw "backend/internal/middleware"
)

type Server struct {
	cfg *appcfg.Config
	htt *http.Server
}

func New(cfg *appcfg.Config) *Server {
	// global logger level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	r := chi.NewRouter()

	// Base middlewares
	r.Use(appmw.RequestID)
	r.Use(appmw.Logger)
	r.Use(appmw.Recover)
	r.Use(appmw.SecurityHeaders)
	r.Use(middleware.Compress(5))

	// Body limit
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	corsMw := cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(corsMw)

	// Rate limiting
	r.Use(httprate.LimitByIP(cfg.RateLimitRequests, cfg.RateLimitWindow))

	// Routes
	r.Get("/healthz", apphandlers.Liveness)
	r.Get("/readyz", apphandlers.Readiness)

	// Serve OpenAPI spec from docs/openapi.yaml
	r.Handle("/openapi.yaml", http.FileServer(http.Dir("./docs")))

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Health endpoint for frontend
		r.Get("/health", apphandlers.Health)

		// CSV data endpoint
		r.Get("/csv", apphandlers.GetCsvData)

		// v1 routes
		r.Route("/v1", func(r chi.Router) {
			r.Get("/ping", apphandlers.Ping)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", apphandlers.Login(cfg))
				r.Group(func(r chi.Router) {
					r.Use(appmw.AuthJWT(cfg.JWTSecret))
					r.Get("/me", apphandlers.Me())
				})
			})
		})
	})

	// Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return &Server{
		cfg: cfg,
		htt: &http.Server{
			Addr:              addr,
			Handler:           bodyLimit(r, cfg.MaxRequestBodySize),
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// bodyLimit returns a handler that limits the size of request bodies.
func bodyLimit(next http.Handler, max int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if max > 0 && r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, max)
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	log.Info().Str("addr", s.htt.Addr).Msg("starting server")
	return s.htt.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down server")
	return s.htt.Shutdown(ctx)
}
