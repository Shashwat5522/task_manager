package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vedologic/task-manager/internal/dto"
	"github.com/vedologic/task-manager/internal/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	log         *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

// Register godoc
// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid register request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		h.log.Error("Registration failed", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, resp)
}

// Login godoc
// @Summary User login
// @Description Login user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid login request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.log.Error("Login failed", zap.Error(err))
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
