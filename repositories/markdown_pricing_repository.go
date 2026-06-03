package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type MarkdownPricingRepository struct{}

func NewMarkdownPricingRepository() *MarkdownPricingRepository {
	return &MarkdownPricingRepository{}
}

func (r *MarkdownPricingRepository) FindByProdukID(idProduk string) (*models.MarkdownPricing, error) {
	var mp models.MarkdownPricing
	result := config.DB.Where("id_produk = ?", idProduk).First(&mp)
	return &mp, result.Error
}

func (r *MarkdownPricingRepository) Save(mp *models.MarkdownPricing) error {
	return config.DB.Save(mp).Error
}

func (r *MarkdownPricingRepository) Create(mp *models.MarkdownPricing) error {
	return config.DB.Create(mp).Error
}

func (r *MarkdownPricingRepository) UpdateNullManual(idProduk string) error {
	return config.DB.Model(&models.MarkdownPricing{}).
		Where("id_produk = ?", idProduk).
		Updates(map[string]interface{}{
			"manual_persen":       nil,
			"manual_aktif_sampai": nil,
		}).Error
}