package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth/handler"
	handler2 "github.com/hardikm9850/GoChat/internal/chat/handler"
	http "github.com/hardikm9850/GoChat/internal/contacts/handler"
	"github.com/hardikm9850/authkit/jwt"
	"github.com/hardikm9850/authkit/middleware"
)

func registerRoutes(
	r *gin.Engine,
	jwtManager *jwt.Manager,
	authHandler *handler.AuthHandler,
	contactsHandler *http.ContactsHandler,
	wsHandler *handler2.WSHandler,
) {

	// Group creates a new router group. We should add all the routes that have common middlewares or the same path prefix.
	// For example, all the routes that use a common middleware for authorization could be grouped.

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GoChat backend running",
			"version": "v1.0",
		})
	})

	// -------- AUTH --------
	auth := r.Group("/auth")
	{
		v1 := auth.Group("/v1")
		{
			v1.POST("/register", authHandler.Register)
			v1.POST("/login", authHandler.Login)
		}
	}

	// -------- WEBSOCKET --------
	ws := r.Group("/ws")
	ws.Use(middleware.JWTAuth(*jwtManager))
	{
		ws.GET("/chat", wsHandler.HandleWebSocket)
	}

	// -------- API --------
	api := r.Group("/api")
	api.Use(middleware.JWTAuth(*jwtManager))
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/me", func(c *gin.Context) {
				userID := c.GetString("userID")
				c.JSON(200, gin.H{"user_id": userID})
			})

			v1.POST("/contacts/sync", contactsHandler.SyncContacts)
		}
	}
}
