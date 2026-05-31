package repositories

import (
	"fmt"
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type StokHistoryRepository struct{}

func NewStokHistoryRepository() *StokHistoryRepository {
	return &StokHistoryRepository{}
}

// Generate ID otomatis — format STK0001
func (r *StokHistoryRepository) GetLastNumber() (int, error) {
	var lastID string
	result := config.DB.Model(&models.StokHistory{}).
		Select("id_history").
		Order("id_history DESC").
		Limit(1).
		Pluck("id_history", &lastID)
	if result.Error != nil || lastID == "" {
		return 0, nil
	}
	var number int
	fmt.Sscanf(lastID, "STK%d", &number)
	return number, nil
}

// Create stok history
func (r *StokHistoryRepository) Create(history *models.StokHistory) error {
	return config.DB.Create(history).Error
}

// FindAll — bisa filter by produk atau by jenis
func (r *StokHistoryRepository) FindAll(idProduk, jenis string) ([]models.StokHistory, error) {
	var histories []models.StokHistory
	query := config.DB.
		Preload("Produk").
		Preload("Produk.Kategori").
		Preload("User")

	if idProduk != "" {
		query = query.Where("id_produk = ?", idProduk)
	}
	if jenis != "" {
		query = query.Where("jenis = ?", jenis)
	}

	result := query.Order("tanggal DESC").Find(&histories)
	return histories, result.Error
}

// FindByProduk — history stok per produk
func (r *StokHistoryRepository) FindByProduk(idProduk string) ([]models.StokHistory, error) {
	return r.FindAll(idProduk, "")
}