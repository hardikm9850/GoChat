package app

import (
	"fmt"
	"github.com/hardikm9850/GoChat/internal/auth/repository"
	mysql "github.com/hardikm9850/GoChat/internal/auth/repository/database"
	"github.com/hardikm9850/GoChat/internal/chat/domain"
	chathandler "github.com/hardikm9850/GoChat/internal/chat/handler"
	"github.com/hardikm9850/GoChat/internal/chat/infrastructure"
	"github.com/hardikm9850/GoChat/internal/chat/usecase"
	"github.com/hardikm9850/GoChat/internal/hub"
	"github.com/hardikm9850/GoChat/internal/infra/db"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth/handler"
	"github.com/hardikm9850/GoChat/internal/auth/repository/memory"
	authservice "github.com/hardikm9850/GoChat/internal/auth/service"
	chatrepodb "github.com/hardikm9850/GoChat/internal/chat/repository/database"
	chatmemory "github.com/hardikm9850/GoChat/internal/chat/repository/memory"
	"github.com/hardikm9850/GoChat/internal/config"
	contacthandler "github.com/hardikm9850/GoChat/internal/contacts/handler"
	contactservice "github.com/hardikm9850/GoChat/internal/contacts/service"
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

	// --- Hub setup
	chatHub := hub.NewHub()
	go chatHub.Run()

	// --- JWT Config ---
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

	// --- Chat ---
	conversationRepo := chatrepodb.New(gormDB)
	conversationRepo.Create(domain.Conversation{
		ID: "conv-123", // arbitrary conversation ID
		Participants: []domain.UserID{
			"df04edf7-8e4a-41f6-a042-7cea242a97e1", // user 1
			"c0b25eb0-2a12-49eb-918a-0f29f12d418a", // user 2
		}, //
	})
	conversationUseCase := usecase.New(conversationRepo)
	messageRepo := chatmemory.NewInMemoryMessageRepository()

	eventPublisher := infrastructure.NewHubEventPublisher(chatHub)

	sendMessageUseCase := usecase.NewSendMessageUseCase(conversationRepo, messageRepo, eventPublisher) // Application layer

	socketHandler := chathandler.NewWSHandler(chatHub, sendMessageUseCase) // Transport layer
	// TODO create conversation handler, repo, use case
	conversationHandler := chathandler.NewConversationHandler(conversationRepo, conversationUseCase)

	registerRoutes(r, &jwtManager, socketHandler, authHandler, contactsHandler, *conversationHandler)

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
