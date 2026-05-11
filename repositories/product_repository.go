package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// FindAll mengambil semua produk aktif dengan pagination
func (r *ProductRepository) FindAll(page, limit int, category string) ([]models.Produk, int64, error) {
	var products []models.Produk
	var total int64

	query := config.DB.Model(&models.Produk{}).Where("is_active = ?", true)

	// Filter by category jika ada
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Hitung total untuk pagination
	query.Count(&total)

	// Ambil data dengan offset & limit
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&products)

	return products, total, result.Error
}

// FindByID mengambil satu produk berdasarkan ID
func (r *ProductRepository) FindByID(id uint) (*models.Produk, error) {
	var product models.Produk
	result := config.DB.First(&product, id)
	return &product, result.Error
}

// Create menyimpan produk baru
func (r *ProductRepository) Create(product *models.Produk) error {
	return config.DB.Create(product).Error
}

// Update memperbarui produk
func (r *ProductRepository) Update(product *models.Produk) error {
	return config.DB.Save(product).Error
}

// Delete soft-delete produk (tidak hapus dari DB)
func (r *ProductRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Produk{}, id).Error
}

// UpdateStock memperbarui stok produk
func (r *ProductRepository) UpdateStock(id uint, stock int) error {
	return config.DB.Model(&models.Produk{}).Where("id = ?", id).Update("stock", stock).Error
}
