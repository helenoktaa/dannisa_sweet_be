// models/markdown_pricing.go
package models

import "time"

type MarkdownPricing struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"     json:"id"`
	IDProduk      string     `gorm:"not null;uniqueIndex;size:20" json:"id_produk"`
	ThresholdHari int        `gorm:"not null;default:2"           json:"threshold_hari"`
	PorsenDiskon  float64    `gorm:"not null;default:20"          json:"persen_diskon"`
	ManualPersen      *float64   `gorm:"default:null"             json:"manual_persen"`
	ManualAktifSampai *time.Time `gorm:"default:null"             json:"manual_aktif_sampai"`
	AktifOtomatis bool       `gorm:"not null;default:true"        json:"aktif_otomatis"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Produk Produk `gorm:"foreignKey:IDProduk;references:IDProduk" json:"produk,omitempty"`
}

type SetMarkdownRequest struct {
	IDProduk      string  `json:"id_produk"      binding:"required"`
	ThresholdHari int     `json:"threshold_hari" binding:"required,min=1"`
	PorsenDiskon  float64 `json:"persen_diskon"  binding:"required,min=0,max=100"`
	AktifOtomatis bool    `json:"aktif_otomatis"`
}

type OverrideMarkdownRequest struct {
	ManualPersen      float64 `json:"manual_persen"       binding:"required,min=0,max=100"`
	ManualAktifSampai string  `json:"manual_aktif_sampai" binding:"required"`
}

type MarkdownPricingResponse struct {
	ID                uint       `json:"id"`
	IDProduk          string     `json:"id_produk"`
	ThresholdHari     int        `json:"threshold_hari"`
	PorsenDiskon      float64    `json:"persen_diskon"`
	ManualPersen      *float64   `json:"manual_persen"`
	ManualAktifSampai *time.Time `json:"manual_aktif_sampai"`
	AktifOtomatis     bool       `json:"aktif_otomatis"`
}

type HargaEfektifResponse struct {
	IDProduk          string   `json:"id_produk"`
	NamaProduk        string   `json:"nama_produk"`
	HargaJual         float64  `json:"harga_jual"`
	PorsenDiskon      float64  `json:"persen_diskon"`
	HargaDiskon       float64  `json:"harga_diskon"`
	SumberDiskon      string   `json:"sumber_diskon"` // "otomatis" / "manual" / "tidak_ada"
	HariMenujuExpired *int     `json:"hari_menuju_expired"`
}