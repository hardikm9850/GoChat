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
	authservice "github.com/hardikm9850/GoChat/internal/auth/service"
	contactservice "github.com/hardikm9850/GoChat/internal/contacts/service"
	contacthandler "github.com/hardikm9850/GoChat/internal/contacts/handler"
	"github.com/hardikm9850/GoChat/internal/config"
	"github.com/hardikm9850/authkit/jwt"
)

type App struct {
	Router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GoChat backend running",
			"version": "v1.0",
		})
	})
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
	// --- Authentication ---
	authRepo, err := buildUserRepository(*cfg, gormDB)
	if err != nil {
		log.Fatal("Failed to setup user repository")
	}
	authService := authservice.New(authRepo, jwtManager)
	authHandler := handler.New(authService)

	// --- Contacts Sync ---
	contactService := contactservice.New(authRepo)
	contactsHandler := contacthandler.NewContactsHandler(contactService)

	registerRoutes(r, &jwtManager, authHandler, contactsHandler)

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
