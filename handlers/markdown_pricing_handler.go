package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type MarkdownPricingHandler struct {
	service *services.MarkdownPricingService
}

func NewMarkdownPricingHandler() *MarkdownPricingHandler {
	return &MarkdownPricingHandler{
		service: services.NewMarkdownPricingService(),
	}
}

// POST /admin/markdown
func (h *MarkdownPricingHandler) SetMarkdown(c *gin.Context) {
	var req models.SetMarkdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	resp, err := h.service.SetMarkdown(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GET /admin/markdown/:id_produk
func (h *MarkdownPricingHandler) GetMarkdown(c *gin.Context) {
	idProduk := c.Param("id_produk")
	resp, err := h.service.GetByProdukID(idProduk)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Konfigurasi markdown tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// PATCH /admin/markdown/:id_produk/override
func (h *MarkdownPricingHandler) OverrideManual(c *gin.Context) {
	idProduk := c.Param("id_produk")
	var req models.OverrideMarkdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	resp, err := h.service.OverrideManual(idProduk, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// DELETE /admin/markdown/:id_produk/override
func (h *MarkdownPricingHandler) HapusOverrideManual(c *gin.Context) {
	idProduk := c.Param("id_produk")
	if err := h.service.HapusOverrideManual(idProduk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Override manual berhasil dihapus"})
}