package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type JabatanHandler struct {
	jabatanService *services.JabatanService
}

func NewJabatanHandler() *JabatanHandler {
	return &JabatanHandler{
		jabatanService: services.NewJabatanService(),
	}
}

// GetAll - GET /v1/jabatan
func (h *JabatanHandler) GetAll(c *gin.Context) {
	jabatans, err := h.jabatanService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data jabatan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jabatans,
	})
}

// GetByID - GET /v1/jabatan/:id
func (h *JabatanHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	jabatan, err := h.jabatanService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Jabatan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jabatan,
	})
}

// Create - POST /v1/jabatan (Admin only)
func (h *JabatanHandler) Create(c *gin.Context) {
	var req models.CreateJabatanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	jabatan, err := h.jabatanService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Jabatan berhasil dibuat",
		"data":    jabatan,
	})
}

// Update - PUT /v1/jabatan/:id (Admin only)
func (h *JabatanHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateJabatanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	jabatan, err := h.jabatanService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jabatan berhasil diperbarui",
		"data":    jabatan,
	})
}

// Delete - DELETE /v1/jabatan/:id (Admin only)
func (h *JabatanHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.jabatanService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jabatan berhasil dihapus",
	})
}