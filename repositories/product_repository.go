package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// FindAll mengambil semua produk
func (r *ProductRepository) FindAll(
	statusProduk string,
) ([]models.Produk, error) {

	var products []models.Produk

	query := config.DB.
		Preload("Kategori").
		Preload("MarkdownPricing")

	if statusProduk != "" {

		query = query.Where(
			"status_produk = ?",
			statusProduk,
		)
	}

	result := query.Find(
		&products,
	)

	return products,
		result.Error
}

// FindByID mengambil produk berdasarkan ID
func (r *ProductRepository) FindByID(id string) (*models.Produk, error) {
	var product models.Produk

	result := config.DB.
		Preload("Kategori").
		Preload("MarkdownPricing").
		Where("id_produk = ?", id).
		First(&product)

	return &product, result.Error
}

// Create menyimpan produk baru
func (r *ProductRepository) Create(product *models.Produk) error {
	return config.DB.Create(product).Error
}

// Update memperbarui data produk
func (r *ProductRepository) Update(product *models.Produk) error {
	return config.DB.Save(product).Error
}

// Delete menghapus produk
func (r *ProductRepository) Delete(id string) error {
	return config.DB.
		Where("id_produk = ?", id).
		Delete(&models.Produk{}).Error
}

// UpdateStock update stok produk
func (r *ProductRepository) UpdateStock(id string, stock int) error {
	return config.DB.
		Model(&models.Produk{}).
		Where("id_produk = ?", id).
		Update("stok", stock).Error
}

// FindNearExpired mengambil produk yang expired dalam N hari ke depan
func (r *ProductRepository) FindNearExpired(hariThreshold int) ([]models.Produk, error) {
    var products []models.Produk

    result := config.DB.
        Preload("Kategori").
        Preload("MarkdownPricing").
        Where("expired_date IS NOT NULL AND expired_date <= NOW() + INTERVAL ? DAY AND expired_date >= NOW()", hariThreshold).
        Find(&products)

    return products, result.Error
}
