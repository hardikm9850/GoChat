package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth/handler"
	http "github.com/hardikm9850/GoChat/internal/contacts/handler"
	"github.com/hardikm9850/authkit/jwt"
	"github.com/hardikm9850/authkit/middleware"
)

func registerRoutes(
	r *gin.Engine,
	jwtManager *jwt.Manager,
	authHandler *handler.AuthHandler,
	contactsHandler *http.ContactsHandler,
) {

	// Group creates a new router group. We should add all the routes that have common middlewares or the same path prefix.
	// For example, all the routes that use a common middleware for authorization could be grouped.
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.JWTAuth(*jwtManager))
	{
		api.GET("/me", func(c *gin.Context) {
			userID := c.GetString("userID")
			c.JSON(200, gin.H{"user_id": userID})
		})
		api.POST("/contacts/sync", contactsHandler.SyncContacts)
	}
}
