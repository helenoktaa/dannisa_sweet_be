package handlers

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type StokHistoryHandler struct {
	service *services.StokHistoryService
}

func NewStokHistoryHandler() *StokHistoryHandler {
	return &StokHistoryHandler{
		service: services.NewStokHistoryService(),
	}
}

// POST /v1/stok-history — catat perubahan stok
func (h *StokHistoryHandler) Create(c *gin.Context) {
	var req models.CreateStokHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Ambil id_user dari JWT
	idUserRaw, _ := c.Get("id_user")
	idUser := fmt.Sprintf("%v", idUserRaw)

	result, err := h.service.Create(req, idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Stok berhasil diperbarui",
		"data":    result,
	})
}

// GET /v1/stok-history — ambil semua history
func (h *StokHistoryHandler) GetAll(c *gin.Context) {
	idProduk := c.Query("id_produk") 
	jenis := c.Query("jenis")        
	tanggalMulai := c.Query("tanggal_mulai")
	tanggalAkhir := c.Query("tanggal_akhir") 

	results, err := h.service.GetAll(idProduk, jenis, tanggalMulai, tanggalAkhir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil history stok",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}