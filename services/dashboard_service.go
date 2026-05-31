package services

import (
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
		Where("status_pembayaran = ?", "Pending").
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