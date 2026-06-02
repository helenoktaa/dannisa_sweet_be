package repositories

import (
	"fmt"

	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"gorm.io/gorm"
)

type TransaksiRepository struct{}

func NewTransaksiRepository() *TransaksiRepository {
	return &TransaksiRepository{}
}

// FindAll mengambil semua transaksi, bisa filter by tanggal
func (r *TransaksiRepository) FindAll(tanggalMulai, tanggalAkhir, status string) ([]models.Transaksi, error) {
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

	if status != "" {
		query = query.Where("status_pembayaran = ?", status)
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

func (r *TransaksiRepository) Create(transaksi *models.Transaksi) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		details := transaksi.Detail
		transaksi.Detail = nil

		// 1. Simpan transaksi header
		if err := tx.Create(transaksi).Error; err != nil {
			return err
		}

		// 2. Simpan detail satu per satu
		for _, detail := range details {
			detail.IDTransaksi = transaksi.IDTransaksi

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			// Kurangi stok HANYA untuk ready stock
			// Pre order stok = 0, tidak perlu dikurangi
			if transaksi.JenisOrder == models.JenisReadyStock {
				if err := tx.Model(&models.Produk{}).
					Where("id_produk = ?", detail.IDProduk).
					UpdateColumn("stok", gorm.Expr("stok - ?", detail.Qty)).
					Error; err != nil {
					return err
				}
			}
		}

		// 3. Kembalikan detail ke struct
		transaksi.Detail = details
		return nil
	})
}

// UpdateStatus update status pembayaran dan jumlah bayar transaksi
func (r *TransaksiRepository) UpdateStatusPembayaran(id string, status string, jumlahBayar float64) error {
	updates := map[string]interface{}{
		"status_pembayaran": status,
	}
	if jumlahBayar > 0 {
		updates["jumlah_bayar"] = jumlahBayar
	}
	return config.DB.Model(&models.Transaksi{}).
		Where("id_transaksi = ?", id).
		Updates(updates).Error
}

// GetLastNumber ambil nomor urut terakhir dari id_transaksi
// Format ID: TDS0001 → ambil angka 0001 → return 1
func (r *TransaksiRepository) GetLastNumber() (int, error) {
	var lastID string
	result := config.DB.Model(&models.Transaksi{}).
		Select("id_transaksi").
		Order("id_transaksi DESC").
		Limit(1).
		Pluck("id_transaksi", &lastID)

	if result.Error != nil || lastID == "" {
		return 0, nil // belum ada transaksi → mulai dari 0
	}

	// Parse angka dari "TDS0001" → 1
	var number int
	fmt.Sscanf(lastID, "TDS%d", &number)
	return number, nil
}

// UpdateStatusOrder - update jenis_order, status_order, catatan
func (r *TransaksiRepository) UpdateStatusOrder(id, statusOrder, catatan string) error {
	updates := map[string]interface{}{
		"status_order": statusOrder,
	}
	if catatan != "" {
		updates["catatan"] = catatan
	}

	result := config.DB.Model(&models.Transaksi{}).
		Where("id_transaksi = ?", id).
		Updates(updates)

	return result.Error
}

// FindPreOrderAktif - semua pre order yang belum selesai/batal
func (r *TransaksiRepository) FindPreOrderAktif() ([]models.Transaksi, error) {
	var transaksis []models.Transaksi

	err := config.DB.
		Where("jenis_order = ? AND status_order NOT IN ?",
			models.JenisPreOrder,
			[]string{models.StatusSelesai, models.StatusDibatalkan},
		).
		Preload("User.Jabatan").
		Preload("Detail.Produk.Kategori").
		Order("tanggal_transaksi DESC").
		Find(&transaksis).Error

	return transaksis, err
}
