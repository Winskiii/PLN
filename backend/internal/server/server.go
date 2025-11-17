package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	appcfg "backend/internal/config"
	"backend/internal/database"
	apphandlers "backend/internal/handlers"
	appmw "backend/internal/middleware"
	"backend/internal/models"
)

type Server struct {
	cfg *appcfg.Config
	htt *http.Server
	db  *sql.DB
}

func New(cfg *appcfg.Config) *Server {
	// global logger level
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Initialize database
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	dbMaxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN", "25"))
	dbMaxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE", "5"))

	dbCfg := &database.DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		Database: getEnv("DB_DATABASE", "work_management"),
		Username: getEnv("DB_USERNAME", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		MaxOpen:  dbMaxOpen,
		MaxIdle:  dbMaxIdle,
	}

	db, err := database.NewMySQL(dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	// Initialize handlers
	authHandler := apphandlers.NewAuthHandler(db, cfg)
	userHandler := apphandlers.NewUserHandler(db, cfg)
	projectHandler := apphandlers.NewProjectHandler(db, cfg)
	taskHandler := apphandlers.NewTaskHandler(db, cfg)

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

			// Auth routes (public)
			r.Route("/auth", func(r chi.Router) {
				r.Post("/login", authHandler.Login)
				r.Group(func(r chi.Router) {
					r.Use(appmw.AuthJWT(cfg.JWTSecret))
					r.Get("/me", authHandler.GetMe)
				})
			})

			// Protected routes (require authentication)
			r.Group(func(r chi.Router) {
				r.Use(appmw.AuthJWT(cfg.JWTSecret))

				// Users management (Admin only)
				r.Route("/users", func(r chi.Router) {
					r.With(appmw.RequireRole(models.RoleAdmin)).Get("/", userHandler.List)
					r.With(appmw.RequireRole(models.RoleAdmin)).Post("/", userHandler.Create)
					r.With(appmw.RequireRole(models.RoleAdmin)).Get("/{id}", userHandler.Get)
					r.With(appmw.RequireRole(models.RoleAdmin)).Put("/{id}", userHandler.Update)
					r.With(appmw.RequireRole(models.RoleAdmin)).Delete("/{id}", userHandler.Delete)
				})

				// Projects management
				r.Route("/projects", func(r chi.Router) {
					r.Get("/", projectHandler.List)
					r.Post("/", projectHandler.Create)
					r.Get("/{id}", projectHandler.Get)
					r.Put("/{id}", projectHandler.Update)
					r.Delete("/{id}", projectHandler.Delete)
				})

				// Tasks management
				r.Route("/tasks", func(r chi.Router) {
					r.Get("/", taskHandler.List)
					r.Get("/my", taskHandler.GetMyTasks)
					r.Post("/", taskHandler.Create)
					r.Get("/{id}", taskHandler.Get)
					r.Put("/{id}", taskHandler.Update)
					r.Delete("/{id}", taskHandler.Delete)
				})
			})
		})
	})

	// Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return &Server{
		cfg: cfg,
		db:  db,
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

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
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
	if s.db != nil {
		s.db.Close()
	}
	return s.htt.Shutdown(ctx)
}
