package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/cleonvoid/course-rating/internal/db"
	"github.com/cleonvoid/course-rating/internal/handler"
	"github.com/cleonvoid/course-rating/internal/middleware"
)

//go:embed internal/templates
var templateFS embed.FS

//go:embed migrations
var migrationsFS embed.FS

func main() {
	var logHandler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL environment variable is required")
		os.Exit(1)
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "dev-secret-change-in-production"
	}

	if err := runMigrations(dbURL); err != nil {
		logger.Error("migration failed", "err", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("database ping failed", "err", err)
		os.Exit(1)
	}

	tmpls, err := parseTemplates()
	if err != nil {
		logger.Error("failed to parse templates", "err", err)
		os.Exit(1)
	}

	queries := db.New(pool)
	courseHandler := handler.NewCourseHandler(queries, tmpls, sessionSecret)
	ratingHandler := handler.NewRatingHandler(queries, courseHandler, sessionSecret)
	enrollHandler := handler.NewEnrollmentHandler(queries, courseHandler, sessionSecret)
	authHandler := handler.NewAuthHandler(queries, tmpls, sessionSecret)
	lessonHandler := handler.NewLessonHandler(queries, tmpls, sessionSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", courseHandler.List)
	mux.HandleFunc("GET /courses/{id}", courseHandler.Detail)
	mux.HandleFunc("POST /courses/{id}/rate", ratingHandler.Upsert)
	mux.HandleFunc("POST /courses/{id}/enroll", enrollHandler.Enroll)
	mux.HandleFunc("GET /courses/{id}/lessons/{lid}", lessonHandler.Detail)
	mux.HandleFunc("POST /courses/{id}/lessons/{lid}/comments", lessonHandler.CreateComment)
	mux.HandleFunc("POST /courses/{id}/lessons/{lid}/comments/{cid}/delete", lessonHandler.DeleteComment)
	mux.HandleFunc("GET /signin", authHandler.SignInPage)
	mux.HandleFunc("POST /signin", authHandler.SignIn)
	mux.HandleFunc("POST /signout", authHandler.SignOut)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rl := middleware.NewIPRateLimiter(20, 50, logger)
	h := middleware.Chain(mux,
		middleware.Logging(logger),
		middleware.Recovery(logger),
		rl.Middleware(),
	)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	go func() {
		logger.Info("listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "err", err)
		os.Exit(1)
	}
	logger.Info("shutdown complete")
}

func runMigrations(dbURL string) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func parseTemplates() (*template.Template, error) {
	funcs := template.FuncMap{
		"starsDisplay": func(avg float64) string {
			filled := int(math.Round(avg))
			var b strings.Builder
			for i := 1; i <= 5; i++ {
				if i <= filled {
					b.WriteRune('★')
				} else {
					b.WriteRune('☆')
				}
			}
			return b.String()
		},
		"seqDesc": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = n - i
			}
			return s
		},
		"dict": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of arguments")
			}
			m := make(map[string]any, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				key, ok := pairs[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: keys must be strings")
				}
				m[key] = pairs[i+1]
			}
			return m, nil
		},
	}

	return template.New("").Funcs(funcs).ParseFS(templateFS,
		"internal/templates/base.html",
		"internal/templates/courses.html",
		"internal/templates/course.html",
		"internal/templates/lesson.html",
		"internal/templates/signin.html",
		"internal/templates/partials/star_rating.html",
		"internal/templates/partials/enroll_section.html",
	)
}
