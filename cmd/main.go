package main

import (
	"log"

	"github.com/hardikm9850/GoChat/internal/app"
	"github.com/hardikm9850/GoChat/internal/config"
)

func main() {
	cfg := config.Load()

	application := app.NewApp(cfg)
	log.Fatal(application.Router.Run(":" + cfg.ServerPort))
}
