package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hardikm9850/GoChat/internal/config"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		log.Printf("DB connect attempt %d\n", i)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("DB connected successfully")
			return db
		}
		log.Println("DB connection failed:", err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("DB connection failed after retries:", err)
	return nil
}
