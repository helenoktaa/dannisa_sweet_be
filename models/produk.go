package models

import "time"

// Produk menyimpan data produk dessert dan kue Dannisa Sweet
type Produk struct {
	IDProduk      string     `gorm:"primaryKey;size:20"                        json:"id_produk"`
	NamaProduk    string     `gorm:"not null;size:100"                         json:"nama_produk"`
	HargaModal    float64    `gorm:"not null"                                  json:"harga_modal"`
	HargaJual     float64    `gorm:"not null"                                  json:"harga_jual"`
	Stok          int        `gorm:"not null;default:0"                        json:"stok"`
	IDKategori    string     `gorm:"not null;index"                            json:"id_kategori"`
	StatusProduk  string     `gorm:"not null;size:20;default:'ready'"          json:"status_produk"` // ready / preorder
	ExpiredDate   *time.Time `gorm:"default:null"                              json:"expired_date"`
	ImageURL     *string    `gorm:"size:500;default:null"                      json:"image_url"`  // nullable

	// Relasi
	Kategori        Kategori          `gorm:"foreignKey:IDKategori;references:IDKategori" json:"kategori,omitempty"`
	DetailTransaksi []DetailTransaksi `gorm:"foreignKey:IDProduk"                         json:"detail_transaksi,omitempty"`
}

// DTO
type CreateProdukRequest struct {
	IDProduk     string  `json:"id_produk"      binding:"required"`
	NamaProduk   string  `json:"nama_produk"    binding:"required"`
	HargaModal   float64 `json:"harga_modal"    binding:"required,min=0"`
	HargaJual    float64 `json:"harga_jual"     binding:"required,min=0"`
	Stok         int     `json:"stok"           binding:"required,min=0"`
	IDKategori   string  `json:"id_kategori"    binding:"required"`
	StatusProduk string  `json:"status_produk"` // ready / preorder, default ready
	ExpiredDate  string  `json:"expired_date"`  // format: 2026-12-31, opsional
	ImageURL     string  `json:"image_url"`
}

type UpdateProdukRequest struct {
	NamaProduk   string  `json:"nama_produk"`
	HargaModal   float64 `json:"harga_modal"  binding:"omitempty,min=0"`
	HargaJual    float64 `json:"harga_jual"   binding:"omitempty,min=0"`
	Stok         int     `json:"stok"         binding:"omitempty,min=0"`
	IDKategori   string  `json:"id_kategori"`
	StatusProduk string  `json:"status_produk"`
	ExpiredDate  string  `json:"expired_date"` // format: 2026-12-31
	ImageURL     string  `json:"image_url"`
}

type UpdateStokRequest struct {
	Stok int `json:"stok" binding:"required,min=0"`
}

// Response
type ProdukResponse struct {
	IDProduk     string           `json:"id_produk"`
	NamaProduk   string           `json:"nama_produk"`
	HargaModal   float64          `json:"harga_modal"`
	HargaJual    float64          `json:"harga_jual"`
	Stok         int              `json:"stok"`
	StatusProduk string           `json:"status_produk"`
	ExpiredDate  *time.Time       `json:"expired_date"`
	ImageURL     *string          `json:"image_url"` 
	Kategori     KategoriResponse `json:"kategori"`
}

type ProdukListResponse struct {
	Data  []ProdukResponse `json:"data"`
	Total int64            `json:"total"`
}