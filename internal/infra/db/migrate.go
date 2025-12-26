package db

import (
	"github.com/hardikm9850/GoChat/internal/auth/domain"
	"github.com/hardikm9850/GoChat/internal/chat/repository/database"
	"gorm.io/gorm"

	authmysql "github.com/hardikm9850/GoChat/internal/auth/repository/database"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&authmysql.UserModel{},
		&database.Conversation{},
		&database.ConversationParticipant{},
		&domain.RefreshToken{},
	)
}
