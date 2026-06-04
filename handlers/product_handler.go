package handlers

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: services.NewProductService(),
	}
}

// GetAll - GET /products
func (h *ProductHandler) GetAll(c *gin.Context) {

	statusProduk := c.Query(
		"status_produk",
	)

	products, err :=
		h.productService.GetAll(
			statusProduk,
		)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,

			gin.H{
				"success": false,
				"message": "Gagal mengambil data produk",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,

		gin.H{
			"success": true,
			"data":    products,
		},
	)
}

// GetByID - GET /products/:id
// GetByID - GET /products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// Sekarang return *ProdukResponse, bukan *Produk
	product, err := h.productService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

// Create - POST /products
// Create - POST /products
func (h *ProductHandler) Create(c *gin.Context) {
    var req models.CreateProdukRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    // Ambil id_user dari JWT context
    idUserRaw, exists := c.Get("id_user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "message": "User tidak terautentikasi",
        })
        return
    }
    idUser := fmt.Sprintf("%v", idUserRaw)

    product, err := h.productService.Create(&req, idUser) // ← tambah idUser
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Produk berhasil dibuat",
        "data":    product,
    })
}

// Update - PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateProdukRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	product, err := h.productService.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk berhasil diperbarui",
		"data":    product,
	})
}

// Delete - DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Produk tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk berhasil dihapus",
	})
}
