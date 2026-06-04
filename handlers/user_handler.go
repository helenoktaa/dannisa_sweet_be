package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type UserHandler struct {
	userService     *services.UserService
	userMenuService *services.UserMenuService
}

func NewUserHandler() *UserHandler {
	// ── Wire dependencies ──────────────────────────────────
	userRepo     := repositories.NewUserRepository()
	userMenuRepo := repositories.NewUserMenuRepository()

	userMenuService := services.NewUserMenuService(userMenuRepo)
	userService     := services.NewUserService(userRepo, userMenuService)

	return &UserHandler{
		userService:     userService,
		userMenuService: userMenuService,
	}
}

// ── GET /v1/users ──────────────────────────────────────────
// Admin only

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// ── GET /v1/users/:id ──────────────────────────────────────
// Admin only

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
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

// ── POST /v1/users ─────────────────────────────────────────
// Admin only

func (h *UserHandler) Create(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := h.userService.Create(req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email sudah terdaftar" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User berhasil ditambahkan",
	})
}

// ── PUT /v1/users/:id ──────────────────────────────────────
// Admin only

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := h.userService.Update(id, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Ambil data terbaru untuk dikembalikan ke client
	user, err := h.userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data user berhasil diperbarui",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data user berhasil diperbarui",
		"data":    user,
	})
}

// ── DELETE /v1/users/:id ───────────────────────────────────
// Admin only

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.Delete(id); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User berhasil dihapus",
	})
}

// ── PUT /v1/users/:id/password ─────────────────────────────
// Admin atau user sendiri

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := h.userService.UpdatePassword(id, req); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = http.StatusNotFound
		}
		if err.Error() == "password lama tidak sesuai" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil diperbarui",
	})
}

// ── GET /v1/users/:id/menus ────────────────────────────────
// Ambil daftar menu yang bisa diakses user

func (h *UserHandler) GetUserMenus(c *gin.Context) {
	id := c.Param("id")

	menuKeys, err := h.userMenuService.GetByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id_user":   id,
			"menu_keys": menuKeys,
		},
	})
}

// ── PUT /v1/users/:id/menus ────────────────────────────────
// Update menu permissions user (Admin only)

func (h *UserHandler) UpdateUserMenus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		MenuKeys []string `json:"menu_keys" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "menu_keys wajib diisi",
		})
		return
	}

	if err := h.userMenuService.Replace(id, req.MenuKeys); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Kembalikan menu terbaru
	menuKeys, _ := h.userMenuService.GetByUserID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Akses menu berhasil diperbarui",
		"data": gin.H{
			"id_user":   id,
			"menu_keys": menuKeys,
		},
	})
}

// ── GET /v1/menus ──────────────────────────────────────────
// Ambil semua menu yang tersedia di aplikasi

func (h *UserHandler) GetAvailableMenus(c *gin.Context) {
	menus := h.userMenuService.GetAllAvailableMenus()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    menus,
	})
}