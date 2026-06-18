package services

import (
	"fmt"
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetDashboard - ambil semua data dashboard sekaligus
func (s *DashboardService) GetDashboard() (*models.DashboardResponse, error) {
	response := &models.DashboardResponse{
		TransaksiPending:       []models.TransaksiPending{},
		ProdukMendekatiExpired: []models.ProdukExpired{},
		ProdukStokMenipis:      []models.ProdukStok{},
	}
	now := time.Now()

	// ── 1. Transaksi Pending ─────────────────────────────────
	// Query semua transaksi yang status_pembayaran = Pending
	// diurutkan dari yang paling lama (prioritas diproses duluan)
	var transaksis []models.Transaksi
	if err := config.DB.
		Where("status_pembayaran IN ? AND status_order != ?", []string{"Pending", "DP"}, "dibatalkan").
		Order("tanggal_transaksi ASC").
		Find(&transaksis).Error; err != nil {
		return nil, err
	}

	response.TotalPending = int64(len(transaksis))
	for _, t := range transaksis {
		hariMenunggu := int(now.Sub(t.TanggalTransaksi).Hours() / 24)
		response.TransaksiPending = append(response.TransaksiPending, models.TransaksiPending{
			IDTransaksi:      t.IDTransaksi,
			NamaCustomer:     t.NamaCustomer,
			JumlahBayar:      t.JumlahBayar,
			MetodePembayaran: t.MetodePembayaran,
			TanggalTransaksi: t.TanggalTransaksi.Format("2006-01-02"),
			HariMenunggu:     hariMenunggu,
			SudahLewat3Hari:  hariMenunggu > 3,
		})
	}

	// ── 2. Produk Mendekati Expired (dalam 7 hari ke depan) ──
	// Hanya produk yang masih ada stoknya
	batas7Hari := now.AddDate(0, 0, 7)
	var produkExpired []models.Produk
	if err := config.DB.
		Where("expired_date IS NOT NULL AND expired_date BETWEEN ? AND ? AND stok > 0",
			now, batas7Hari).
		Order("expired_date ASC").
		Find(&produkExpired).Error; err != nil {
		return nil, err
	}

	response.TotalMendekatiExpired = int64(len(produkExpired))
	for _, p := range produkExpired {
		if p.ExpiredDate == nil {
			continue
		}
		sisaHari := int(p.ExpiredDate.Sub(now).Hours() / 24)
		response.ProdukMendekatiExpired = append(response.ProdukMendekatiExpired, models.ProdukExpired{
			IDProduk:    p.IDProduk,
			NamaProduk:  p.NamaProduk,
			Stok:        p.Stok,
			ExpiredDate: p.ExpiredDate.Format("2006-01-02"),
			SisaHari:    sisaHari,
		})
	}

	// ── 3. Produk Stok Menipis (stok <= 5) ──────────────────
	// Diurutkan dari stok paling sedikit
	var produkMenipis []models.Produk
	if err := config.DB.
		Where("stok > 0 AND stok <= 5").
		Order("stok ASC").
		Find(&produkMenipis).Error; err != nil {
		return nil, err
	}

	response.TotalStokMenipis = int64(len(produkMenipis))
	for _, p := range produkMenipis {
		response.ProdukStokMenipis = append(response.ProdukStokMenipis, models.ProdukStok{
			IDProduk:     p.IDProduk,
			NamaProduk:   p.NamaProduk,
			Stok:         p.Stok,
			StatusProduk: p.StatusProduk,
		})
	}

	return response, nil
}

func (s *DashboardService) GetDashboardHarian() (*models.DashboardHarian, error) {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(wib)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, wib)
	todayEnd := todayStart.Add(24 * time.Hour)
	threeDaysAgo := todayStart.AddDate(0, 0, -3)

	fmt.Printf("now WIB: %v\n", now)
	fmt.Printf("todayStart: %v\n", todayStart)
	fmt.Printf("todayEnd: %v\n", todayEnd)

	result := &models.DashboardHarian{
		TransaksiTerbaru: []models.TransaksiTerbaru{},
	}

	// 1. Total pending > 3 hari
	if err := config.DB.Model(&models.Transaksi{}).
		Where("status_pembayaran IN ? AND tanggal_transaksi < ? AND status_order != ?", []string{"Pending", "DP"}, threeDaysAgo, "dibatalkan").
		Count(&result.TotalPendingLewat3Hari).Error; err != nil {
		return nil, err
	}

	// 2. Total lunas hari ini
	if err := config.DB.Model(&models.Transaksi{}).
		Where("status_pembayaran = ? AND tanggal_lunas >= ? AND tanggal_lunas < ?",
			"Lunas", todayStart, todayEnd).
		Count(&result.TotalLunasHariIni).Error; err != nil {
		return nil, err
	}

	// 3. Omzet & total transaksi — hitung dari detail, bukan jumlah_bayar
	type OmzetResult struct {
		TotalOmzet     float64 `gorm:"column:total_omzet"`
		TotalTransaksi int64   `gorm:"column:total_transaksi"`
	}
	var omzet OmzetResult
	if err := config.DB.Table("detail_transaksi dt").
		Select(`COALESCE(SUM(dt.qty * dt.harga_jual), 0) as total_omzet,
        COUNT(DISTINCT t.id_transaksi) as total_transaksi`).
		Joins("JOIN transaksi t ON t.id_transaksi = dt.id_transaksi").
		Where("t.status_pembayaran = ? AND t.tanggal_lunas >= ? AND t.tanggal_lunas < ?",
			"Lunas", todayStart, todayEnd).
		Scan(&omzet).Error; err != nil {
		return nil, err
	}
	result.TotalOmzet = omzet.TotalOmzet
	result.TotalTransaksi = omzet.TotalTransaksi

	// 4. Total modal
	type ModalResult struct {
		TotalModal float64 `gorm:"column:total_modal"`
	}
	var modal ModalResult
	if err := config.DB.Table("detail_transaksi dt").
		Select("COALESCE(SUM(dt.qty * p.harga_modal), 0) as total_modal").
		Joins("JOIN transaksi t ON t.id_transaksi = dt.id_transaksi").
		Joins("JOIN produk p ON p.id_produk = dt.id_produk").
		Where("t.status_pembayaran = ? AND t.tanggal_lunas >= ? AND t.tanggal_lunas < ?",
			"Lunas", todayStart, todayEnd).
		Scan(&modal).Error; err != nil {
		return nil, err
	}
	result.TotalModal = modal.TotalModal
	result.KeuntunganBersih = result.TotalOmzet - result.TotalModal
	// 5. Transaksi terbaru — tambah total_penjualan dari detail
	type TrxRaw struct {
		IDTransaksi      string
		NamaCustomer     string
		TanggalTransaksi time.Time
		TotalItem        int
		TotalPenjualan   float64 // ← ganti dari JumlahBayar
		MetodePembayaran string
		StatusPembayaran string
	}
	var rows []TrxRaw

	if err := config.DB.Table("transaksi t").
		Select(`t.id_transaksi, t.nama_customer, t.tanggal_transaksi,
        COALESCE(SUM(dt.qty), 0) as total_item,
        COALESCE(SUM(dt.qty * dt.harga_jual), 0) as total_penjualan,
        t.metode_pembayaran, t.status_pembayaran`).
		Joins("LEFT JOIN detail_transaksi dt ON dt.id_transaksi = t.id_transaksi").
		Where("t.tanggal_transaksi >= ? AND t.status_order != ?", threeDaysAgo, "dibatalkan").
		Group("t.id_transaksi, t.nama_customer, t.tanggal_transaksi, t.metode_pembayaran, t.status_pembayaran").
		Order("t.tanggal_transaksi DESC").
		Limit(20).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		result.TransaksiTerbaru = append(result.TransaksiTerbaru, models.TransaksiTerbaru{
			IDTransaksi:      r.IDTransaksi,
			NamaCustomer:     r.NamaCustomer,
			TanggalTransaksi: r.TanggalTransaksi.Format("02/01/2006 15:04"),
			TotalItem:        r.TotalItem,
			JumlahBayar:      r.TotalPenjualan, // ← pakai TotalPenjualan
			MetodePembayaran: r.MetodePembayaran,
			StatusPembayaran: r.StatusPembayaran,
		})
	}
	return result, nil
}
