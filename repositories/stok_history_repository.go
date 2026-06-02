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

func (r *StokHistoryRepository) Create(history *models.StokHistory) error {
	return config.DB.Create(history).Error
}

// FindAll — filter by produk, jenis, dan tanggal
func (r *StokHistoryRepository) FindAll(idProduk, jenis, tanggalMulai, tanggalAkhir string) ([]models.StokHistory, error) {
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
	if tanggalMulai != "" {
		query = query.Where("DATE(tanggal) >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("DATE(tanggal) <= ?", tanggalAkhir)
	}

	result := query.Order("tanggal DESC").Find(&histories)
	return histories, result.Error
}

func (r *StokHistoryRepository) FindByProduk(idProduk string) ([]models.StokHistory, error) {
	return r.FindAll(idProduk, "", "", "")
}