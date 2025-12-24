package main

import (
	"fmt"
	"os"

	"github.com/hardikm9850/GoChat/internal/app"
	"github.com/hardikm9850/GoChat/internal/config"
)

func main() {
	cfg := config.Load()

	a := app.NewApp(cfg)

	fmt.Println("Starting server on port:", cfg.ServerPort)
	err := a.Router.Run(":" + cfg.ServerPort)
	if err != nil {
		fmt.Println("Router failed to start:", err)
		os.Exit(1)
	}
}
