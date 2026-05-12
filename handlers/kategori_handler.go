package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type KategoriHandler struct {
	kategoriService *services.KategoriService
}

func NewKategoriHandler() *KategoriHandler {
	return &KategoriHandler{
		kategoriService: services.NewKategoriService(),
	}
}

// GetAll - GET /v1/kategori
func (h *KategoriHandler) GetAll(c *gin.Context) {
	kategoris, err := h.kategoriService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kategoris,
	})
}

// GetByID - GET /v1/kategori/:id
func (h *KategoriHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	kategori, err := h.kategoriService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kategori,
	})
}

// Create - POST /v1/kategori (Admin only)
func (h *KategoriHandler) Create(c *gin.Context) {
	var req models.CreateKategoriRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	kategori, err := h.kategoriService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kategori berhasil dibuat",
		"data":    kategori,
	})
}

// Update - PUT /v1/kategori/:id (Admin only)
func (h *KategoriHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateKategoriRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	kategori, err := h.kategoriService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil diperbarui",
		"data":    kategori,
	})
}

// Delete - DELETE /v1/kategori/:id (Admin only)
func (h *KategoriHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.kategoriService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil dihapus",
	})
}