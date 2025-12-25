package handler

import (
	"errors"
	response "github.com/hardikm9850/GoChat/internal/http/validation"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hardikm9850/GoChat/internal/auth"
	"github.com/hardikm9850/GoChat/internal/auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func New(auth service.AuthService) *AuthHandler {
	return &AuthHandler{service: auth}
}

type registerRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type loginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

// @BasePath /api/v1

// @Summary Register a new user
// @Description Registers a new user with phone, password, and name
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body registerRequest true "Register Request"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Registration failed"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log.Println("Http 👉 /auth/register hit")
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(req.Phone, req.Password, req.Name); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// @Summary Login user
// @Description Login user with phone and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "Login Request"
// @Success 200 {object} loginResponse "Access token returned"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Login failed"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ValidationError(err))
		return
	}

	tokens, err := h.service.Login(req.Phone, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken: tokens.AccessToken,
	})
}
