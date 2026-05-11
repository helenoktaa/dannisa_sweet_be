package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: services.NewAuthService()}
}

// Login godoc
// POST /v1/auth/login
// Terima email + password → verifikasi → return JWT
func (h *AuthHandler) Login(c *gin.Context) {
	// 1. Parse request body
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email dan password wajib diisi",
		})
		return
	}

	// 2. Proses login via service
	loginResp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Return JWT + data user
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data":    loginResp,
	})
}

// Register godoc
// POST /v1/auth/register
// Hanya bisa diakses Admin (tambah akun kasir baru)
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Akun berhasil dibuat",
		"data":    user,
	})
}

// GetProfile godoc
// GET /v1/auth/profile
// Ambil data profil user yang sedang login
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// IDUser diambil dari JWT via middleware
	idUser := c.GetString("id_user")

	user, err := h.authService.GetProfile(idUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateProfile godoc
// PUT /v1/auth/profile
// Update data profil user yang sedang login
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	idUser := c.GetString("id_user")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user, err := h.authService.UpdateProfile(idUser, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profil berhasil diupdate",
		"data":    user,
	})
}

// UpdatePassword godoc
// PUT /v1/auth/password
// Ganti password user yang sedang login
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	idUser := c.GetString("id_user")

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := h.authService.UpdatePassword(idUser, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil diubah",
	})
}