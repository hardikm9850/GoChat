package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/hardikm9850/GoChat/docs"
	"github.com/hardikm9850/GoChat/internal/app"
	"github.com/hardikm9850/GoChat/internal/config"
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title GoChat API
// @version 1.0
// @description Chat backend service
// @host localhost:8080
// @BasePath /api/v1
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("MAIN STARTED")
	log.Println("PORT =", os.Getenv("PORT"))
	log.Println("ENV =", os.Environ())

	cfg := config.Load()
	a := app.NewApp(cfg)

	fmt.Println("Starting server on port:", cfg.ServerPort)
	err := a.Router.Run(":" + cfg.ServerPort)
	if err != nil {
		fmt.Println("Router failed to start:", err)
		os.Exit(1)
	}
}
