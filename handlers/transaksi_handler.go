package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type TransaksiHandler struct {
	transaksiService *services.TransaksiService
}

func NewTransaksiHandler() *TransaksiHandler {
	return &TransaksiHandler{
		transaksiService: services.NewTransaksiService(),
	}
}

// Create - POST /v1/transaksi
// Kasir input transaksi baru beserta detail itemnya
func (h *TransaksiHandler) Create(c *gin.Context) {
	var req models.CreateTransaksiRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Ambil id_user dari JWT (kasir yang sedang login)
	idUser := c.GetString("id_user")
	req.IDUser = idUser

	transaksi, err := h.transaksiService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Transaksi berhasil dibuat",
		"data":    transaksi,
	})
}

// GetAll - GET /v1/transaksi (Admin only)
// Laporan semua transaksi, bisa filter by tanggal
func (h *TransaksiHandler) GetAll(c *gin.Context) {
	// Query params opsional untuk filter tanggal
	tanggalMulai := c.Query("tanggal_mulai")
	tanggalAkhir := c.Query("tanggal_akhir")

	transaksis, err := h.transaksiService.GetAll(tanggalMulai, tanggalAkhir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data transaksi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksis,
	})
}

// GetByID - GET /v1/transaksi/:id
// Detail satu transaksi beserta semua item
func (h *TransaksiHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	transaksi, err := h.transaksiService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Transaksi tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksi,
	})
}

// GetStruk - GET /v1/transaksi/:id/struk
// Generate data invoice dari transaksi
func (h *TransaksiHandler) GetInvoice(c *gin.Context) {
    id := c.Param("id")

    invoice, err := h.transaksiService.GetInvoice(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Transaksi tidak ditemukan",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    invoice,
    })
}
// GetLaporan - GET /v1/transaksi/laporan (Admin only)
// Laporan penjualan dengan total modal, penjualan, dan laba
func (h *TransaksiHandler) GetLaporan(c *gin.Context) {
	var req models.LaporanRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "tanggal_mulai dan tanggal_akhir wajib diisi (format: 2024-01-01)",
		})
		return
	}

	laporan, err := h.transaksiService.GetLaporan(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    laporan,
	})
}

// UpdateStatus - PUT /v1/transaksi/:id/status (Admin only)
// Update status pembayaran: Pending -> Lunas
func (h *TransaksiHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateStatusPembayaranRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	transaksi, err := h.transaksiService.UpdateStatus(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status pembayaran berhasil diperbarui",
		"data":    transaksi,
	})
}