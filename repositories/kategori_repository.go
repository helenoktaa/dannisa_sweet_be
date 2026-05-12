package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type KategoriRepository struct{}

func NewKategoriRepository() *KategoriRepository {
	return &KategoriRepository{}
}

func (r *KategoriRepository) FindAll() ([]models.Kategori, error) {
	var kategoris []models.Kategori
	result := config.DB.Find(&kategoris)
	return kategoris, result.Error
}

func (r *KategoriRepository) FindByID(id string) (*models.Kategori, error) {
	var kategori models.Kategori
	result := config.DB.Where("id_kategori = ?", id).First(&kategori)
	return &kategori, result.Error
}

func (r *KategoriRepository) Create(kategori *models.Kategori) error {
	return config.DB.Create(kategori).Error
}

func (r *KategoriRepository) Update(kategori *models.Kategori) error {
	return config.DB.Save(kategori).Error
}

func (r *KategoriRepository) Delete(id string) error {
	return config.DB.Where("id_kategori = ?", id).Delete(&models.Kategori{}).Error
}