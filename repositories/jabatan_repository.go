package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type JabatanRepository struct{}

func NewJabatanRepository() *JabatanRepository {
	return &JabatanRepository{}
}

func (r *JabatanRepository) FindAll() ([]models.Jabatan, error) {
	var jabatans []models.Jabatan
	result := config.DB.Find(&jabatans)
	return jabatans, result.Error
}

func (r *JabatanRepository) FindByID(id string) (*models.Jabatan, error) {
	var jabatan models.Jabatan
	result := config.DB.Where("id_jabatan = ?", id).First(&jabatan)
	return &jabatan, result.Error
}

func (r *JabatanRepository) Create(jabatan *models.Jabatan) error {
	return config.DB.Create(jabatan).Error
}

func (r *JabatanRepository) Update(jabatan *models.Jabatan) error {
	return config.DB.Save(jabatan).Error
}

func (r *JabatanRepository) Delete(id string) error {
	return config.DB.Where("id_jabatan = ?", id).Delete(&models.Jabatan{}).Error
}