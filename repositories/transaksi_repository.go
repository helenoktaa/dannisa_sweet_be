package repositories

import (
	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"gorm.io/gorm"
)

type TransaksiRepository struct{}

func NewTransaksiRepository() *TransaksiRepository {
	return &TransaksiRepository{}
}

// FindAll mengambil semua transaksi, bisa filter by tanggal
func (r *TransaksiRepository) FindAll(tanggalMulai, tanggalAkhir string) ([]models.Transaksi, error) {
	var transaksis []models.Transaksi

	query := config.DB.
		Preload("User").
		Preload("User.Jabatan").
		Preload("Detail").
		Preload("Detail.Produk").
		Preload("Detail.Produk.Kategori")

	// Filter tanggal kalau ada
	if tanggalMulai != "" && tanggalAkhir != "" {
		query = query.Where(
			"tanggal_transaksi BETWEEN ? AND ?",
			tanggalMulai+" 00:00:00",
			tanggalAkhir+" 23:59:59",
		)
	}

	result := query.Order("tanggal_transaksi DESC").Find(&transaksis)
	return transaksis, result.Error
}

// FindByID mengambil satu transaksi beserta semua relasinya
func (r *TransaksiRepository) FindByID(id string) (*models.Transaksi, error) {
	var transaksi models.Transaksi

	result := config.DB.
		Preload("User").
		Preload("User.Jabatan").
		Preload("Detail").
		Preload("Detail.Produk").
		Preload("Detail.Produk.Kategori").
		Where("id_transaksi = ?", id).
		First(&transaksi)

	return &transaksi, result.Error
}

// Create menyimpan transaksi baru beserta detail-nya dalam satu transaksi DB
func (r *TransaksiRepository) Create(transaksi *models.Transaksi) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Pisahkan detail dari transaksi sebelum insert header
		details := transaksi.Detail
		transaksi.Detail = nil // ← kosongkan dulu biar GORM tidak auto-insert detail

		// 1. Simpan transaksi header dulu
		if err := tx.Create(transaksi).Error; err != nil {
			return err
		}

		// 2. Simpan detail satu per satu + update stok
		for _, detail := range details {
			detail.IDTransaksi = transaksi.IDTransaksi

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			// Kurangi stok produk
			if err := tx.Model(&models.Produk{}).
				Where("id_produk = ?", detail.IDProduk).
				UpdateColumn("stok", gorm.Expr("stok - ?", detail.Qty)).
				Error; err != nil {
				return err
			}
		}

		// 3. Kembalikan detail ke struct (untuk response)
		transaksi.Detail = details
		return nil
	})
}

// UpdateStatus update status pembayaran transaksi
func (r *TransaksiRepository) UpdateStatus(id string, status string) error {
	return config.DB.Model(&models.Transaksi{}).
		Where("id_transaksi = ?", id).
		Update("status_pembayaran", status).Error
}
