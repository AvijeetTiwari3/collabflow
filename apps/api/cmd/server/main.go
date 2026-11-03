package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/sde/collabflow/internal/auth"
	"github.com/sde/collabflow/internal/config"
	"github.com/sde/collabflow/internal/handler"
	"github.com/sde/collabflow/internal/middleware"
	"github.com/sde/collabflow/internal/repository"
	"github.com/sde/collabflow/internal/service"
	"github.com/sde/collabflow/internal/ws"
)

func main() {
	cfg := config.Load()
	log.Printf("[main] starting CollabFlow API Engine (env=%s)...", cfg.Environment)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 1. Connect to PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[main] failed to connect to database: %v", err)
	}
	defer db.Close()

	// Tune connection pool
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Auto-run schema migrations
	if err := runMigrations(db); err != nil {
		log.Printf("[main] warning: migration verification: %v", err)
	} else {
		log.Println("[main] database schema migrations verified successfully")
	}

	// 2. Connect to Redis (optional with graceful fallback)
	var rdb *redis.Client
	if cfg.RedisURL != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:        cfg.RedisURL,
			DialTimeout: 2 * time.Second,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("[main] warning: redis ping failed, operating in single-node in-memory mode: %v", err)
			rdb = nil
		} else {
			log.Println("[main] connected to Redis Pub/Sub cluster")
		}
		cancel()
	}

	// 3. Initialize Token Manager
	tokenMgr := auth.NewTokenManager(cfg.JWTSecret, 15, 7)

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// 5. Initialize Services
	authSvc := service.NewAuthService(userRepo, refreshRepo, tokenMgr)
	projectSvc := service.NewProjectService(projectRepo)
	taskSvc := service.NewTaskService(taskRepo)

	// 6. Initialize WebSocket Hub
	hub := ws.NewHub(rdb)
	go hub.Run()

	// 7. Initialize Handlers
	authH := handler.NewAuthHandler(authSvc)
	projectH := handler.NewProjectHandler(projectSvc)
	taskH := handler.NewTaskHandler(taskSvc)
	searchH := handler.NewSearchHandler(taskSvc)
	wsH := handler.NewWebSocketHandler(hub, tokenMgr)
	healthH := handler.NewHealthHandler(db)

	// 8. Setup Router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.CORSOrigin))

	// Health and Discovery
	router.GET("/", healthH.Root)
	router.HEAD("/", healthH.Root)
	router.GET("/health", healthH.Liveness)
	router.HEAD("/health", healthH.Liveness)
	router.GET("/ready", healthH.Readiness)
	router.HEAD("/ready", healthH.Readiness)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// WebSocket endpoint (handles its own token auth via query param)
	router.GET("/ws", wsH.HandleWS)

	// Public Auth Group
	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", authH.Register)
		v1.POST("/auth/login", authH.Login)
		v1.POST("/auth/refresh", authH.Refresh)
	}

	// Protected Routes (JWT required)
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired(tokenMgr))
	{
		protected.POST("/auth/logout", authH.Logout)
		protected.GET("/auth/me", authH.Me)

		protected.POST("/projects", projectH.Create)
		protected.GET("/projects", projectH.List)
		protected.GET("/projects/:id", projectH.GetByID)

		protected.POST("/projects/:id/tasks", taskH.Create)
		protected.GET("/projects/:id/tasks", taskH.ListByProject)
		protected.PATCH("/tasks/:id", taskH.Update)
		protected.POST("/tasks/:id/reorder", taskH.Reorder) // O(1) Fractional Indexing
		protected.DELETE("/tasks/:id", taskH.Delete)

		protected.GET("/search", searchH.Search)
		protected.GET("/rooms/:id/presence", wsH.GetPresence)
	}

	// 9. Dual IPv4 Server Listeners (Binding explicitly to 0.0.0.0)
	primaryAddr := "0.0.0.0:" + cfg.Port
	srv := &http.Server{
		Addr:         primaryAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("[main] server listening on primary IPv4 %s", primaryAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[main] primary listener failure: %v", err)
		}
	}()

	// Fallback port 10000 for cloud proxy defaults (Render, etc.)
	if cfg.Port != "10000" {
		srv10k := &http.Server{
			Addr:         "0.0.0.0:10000",
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		go func() {
			log.Println("[main] server also listening on fallback IPv4 port 0.0.0.0:10000")
			_ = srv10k.ListenAndServe()
		}()
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[main] graceful shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[main] server forced shutdown: %v", err)
	}

	log.Println("[main] server cleanly stopped")
}

func runMigrations(db *sql.DB) error {
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS users (
		id            UUID PRIMARY KEY,
		email         VARCHAR(255) UNIQUE NOT NULL,
		name          VARCHAR(100) NOT NULL,
		avatar_url    TEXT,
		password_hash VARCHAR(255) NOT NULL,
		created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id          UUID PRIMARY KEY,
		user_id     UUID NOT NULL,
		token_hash  VARCHAR(64) UNIQUE NOT NULL,
		expires_at  TIMESTAMPTZ NOT NULL,
		revoked     BOOLEAN NOT NULL DEFAULT FALSE,
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash) WHERE revoked = FALSE;

	CREATE TABLE IF NOT EXISTS projects (
		id           UUID PRIMARY KEY,
		workspace_id UUID NOT NULL,
		name         VARCHAR(200) NOT NULL,
		description  TEXT,
		created_by   UUID NOT NULL,
		created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id          UUID PRIMARY KEY,
		project_id  UUID NOT NULL,
		title       VARCHAR(500) NOT NULL,
		description TEXT,
		status      VARCHAR(50) NOT NULL DEFAULT 'todo',
		priority    VARCHAR(20) NOT NULL DEFAULT 'medium',
		position    FLOAT NOT NULL DEFAULT 1000.0,
		assignee_id UUID,
		due_date    TIMESTAMPTZ,
		created_by  UUID NOT NULL,
		created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_project_pos ON tasks(project_id, position ASC);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(project_id, status);
	`
	_, err := db.Exec(schema)
	return err
}
