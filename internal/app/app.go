package app

import (
	"fmt"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	mysql "github.com/hardikm9850/GoChat/internal/auth/repository/database"
	"github.com/hardikm9850/GoChat/internal/infra/db"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth/handler"
	"github.com/hardikm9850/GoChat/internal/auth/repository/memory"
	"github.com/hardikm9850/GoChat/internal/auth/service"
	"github.com/hardikm9850/GoChat/internal/config"
	"github.com/hardikm9850/authkit/jwt"
)

type App struct {
	Router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	r := gin.Default()

	// --- Database setup ---
	gormDB := db.Connect(cfg)

	if err := db.Migrate(gormDB); err != nil {
    	log.Fatal("DB migration failed:", err)
	}

	c := jwt.Config{
		Algorithm:      jwt.HS256,
		Secret:         cfg.JWTSecret,
		Issuer:         "authkit",
		Audience:       "authkit",
		AccessTokenTTL: time.Hour * 24 * 7,
	}

	// --- Auth wiring ---
	jwtManager, err := jwt.NewManager(c)
	if err != nil {
		log.Fatal(err)
	}

	authRepo, err := buildUserRepository(*cfg, gormDB)
	if err != nil {
		log.Fatal("Failed to setup user repository")
	}
	authService := service.New(authRepo, jwtManager)
	authHandler := handler.New(authService)

	registerRoutes(r, authHandler, &jwtManager)

	return &App{
		Router: r,
	}
}

func buildUserRepository(cfg config.Config, db *gorm.DB) (repository.UserRepository, error) {
	switch cfg.UserRepoType {
	case config.InMemory:
		return memory.New(), nil

	case config.MySQL:
		return mysql.New(db), nil

	default:
		return nil, fmt.Errorf("invalid USER_REPO: %s", cfg.UserRepoType)
	}
}
