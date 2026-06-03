package services

import (
	"errors"
	"math"
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
	"gorm.io/gorm"
)

type MarkdownPricingService struct {
	repo *repositories.MarkdownPricingRepository
}

func NewMarkdownPricingService() *MarkdownPricingService {
	return &MarkdownPricingService{
		repo: repositories.NewMarkdownPricingRepository(),
	}
}

func (s *MarkdownPricingService) SetMarkdown(req models.SetMarkdownRequest) (*models.MarkdownPricingResponse, error) {
	existing, err := s.repo.FindByProdukID(req.IDProduk)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		mp := &models.MarkdownPricing{
			IDProduk:      req.IDProduk,
			ThresholdHari: req.ThresholdHari,
			PorsenDiskon:  req.PorsenDiskon,
			AktifOtomatis: req.AktifOtomatis,
		}
		if err := s.repo.Create(mp); err != nil {
			return nil, err
		}
		return toMarkdownResponse(mp), nil
	} else if err != nil {
		return nil, err
	}

	existing.ThresholdHari = req.ThresholdHari
	existing.PorsenDiskon = req.PorsenDiskon
	existing.AktifOtomatis = req.AktifOtomatis
	if err := s.repo.Save(existing); err != nil {
		return nil, err
	}
	return toMarkdownResponse(existing), nil
}

func (s *MarkdownPricingService) OverrideManual(idProduk string, req models.OverrideMarkdownRequest) (*models.MarkdownPricingResponse, error) {
	aktifSampai, err := time.Parse("2006-01-02", req.ManualAktifSampai)
	if err != nil {
		return nil, errors.New("format manual_aktif_sampai tidak valid, gunakan YYYY-MM-DD")
	}
	if aktifSampai.Before(time.Now()) {
		return nil, errors.New("manual_aktif_sampai harus di masa depan")
	}

	mp, err := s.repo.FindByProdukID(idProduk)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("konfigurasi markdown belum ada, buat dulu via SET")
	} else if err != nil {
		return nil, err
	}

	mp.ManualPersen = &req.ManualPersen
	mp.ManualAktifSampai = &aktifSampai
	if err := s.repo.Save(mp); err != nil {
		return nil, err
	}
	return toMarkdownResponse(mp), nil
}

func (s *MarkdownPricingService) HapusOverrideManual(idProduk string) error {
	return s.repo.UpdateNullManual(idProduk)
}

func (s *MarkdownPricingService) GetByProdukID(idProduk string) (*models.MarkdownPricingResponse, error) {
	mp, err := s.repo.FindByProdukID(idProduk)
	if err != nil {
		return nil, err
	}
	return toMarkdownResponse(mp), nil
}

func (s *MarkdownPricingService) HitungHargaEfektif(produk models.Produk) models.HargaEfektifResponse {
	resp := models.HargaEfektifResponse{
		IDProduk:     produk.IDProduk,
		NamaProduk:   produk.NamaProduk,
		HargaJual:    produk.HargaJual,
		HargaDiskon:  produk.HargaJual,
		SumberDiskon: "tidak_ada",
	}

	// Hitung sisa hari menuju expired
	var hariExpired *int
	if produk.ExpiredDate != nil {
		sisa := int(math.Ceil(time.Until(*produk.ExpiredDate).Hours() / 24))
		hariExpired = &sisa
	}
	resp.HariMenujuExpired = hariExpired

	// Ambil konfigurasi markdown — pakai yang sudah di-preload jika ada
	var mp *models.MarkdownPricing
	if produk.MarkdownPricing != nil {
		mp = produk.MarkdownPricing
	} else {
		// fallback ke DB jika tidak di-preload
		fetched, err := s.repo.FindByProdukID(produk.IDProduk)
		if err != nil {
			return resp
		}
		mp = fetched
	}

	// Prioritas 1: Override manual
	if mp.ManualPersen != nil && mp.ManualAktifSampai != nil {
		if time.Now().Before(*mp.ManualAktifSampai) {
			resp.PorsenDiskon = *mp.ManualPersen
			resp.HargaDiskon = bulatkan(produk.HargaJual * (1 - *mp.ManualPersen/100))
			resp.SumberDiskon = "manual"
			return resp
		}
		// Expired → bersihkan di goroutine
		go s.HapusOverrideManual(produk.IDProduk)
	}

	// Prioritas 2: Diskon otomatis
	if mp.AktifOtomatis && hariExpired != nil && *hariExpired >= 0 {
		if *hariExpired <= mp.ThresholdHari {
			resp.PorsenDiskon = mp.PorsenDiskon
			resp.HargaDiskon = bulatkan(produk.HargaJual * (1 - mp.PorsenDiskon/100))
			resp.SumberDiskon = "otomatis"
		}
	}

	return resp
}

// --- helpers ---

func bulatkan(harga float64) float64 {
	return math.Round(harga/100) * 100
}

func toMarkdownResponse(mp *models.MarkdownPricing) *models.MarkdownPricingResponse {
	return &models.MarkdownPricingResponse{
		ID:                mp.ID,
		IDProduk:          mp.IDProduk,
		ThresholdHari:     mp.ThresholdHari,
		PorsenDiskon:      mp.PorsenDiskon,
		ManualPersen:      mp.ManualPersen,
		ManualAktifSampai: mp.ManualAktifSampai,
		AktifOtomatis:     mp.AktifOtomatis,
	}
}