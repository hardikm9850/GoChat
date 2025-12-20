package db

import (
    "gorm.io/gorm"

    authmysql "github.com/hardikm9850/GoChat/internal/auth/repository/database"
)

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &authmysql.UserModel{},
    )
}
