package models

import "time"

// StokHistory menyimpan riwayat perubahan stok produk
type StokHistory struct {
	IDHistory   string    `gorm:"primaryKey;size:30"                    json:"id_history"`
	IDProduk    string    `gorm:"not null;size:20;index"                json:"id_produk"`
	IDUser      string    `gorm:"not null;size:20;index"                json:"id_user"`
	Jenis       string    `gorm:"not null;size:20"                      json:"jenis"`       // penambahan / pengurangan
	Jumlah      int       `gorm:"not null"                              json:"jumlah"`
	StokSebelum int       `gorm:"not null"                              json:"stok_sebelum"`
	StokSesudah int       `gorm:"not null"                              json:"stok_sesudah"`
	Keterangan  string    `gorm:"size:255"                              json:"keterangan"`  // rusak / expired / restock / dll
	NilaiRugi   float64   `gorm:"not null;default:0"                 json:"nilai_rugi"`
	Tanggal     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"    json:"tanggal"`

	// Relasi
	Produk Produk `gorm:"foreignKey:IDProduk;references:IDProduk" json:"produk,omitempty"`
	User   User   `gorm:"foreignKey:IDUser;references:IDUser"     json:"user,omitempty"`
}

// DTO
type CreateStokHistoryRequest struct {
	IDProduk   string `json:"id_produk"   binding:"required"`
	Jenis      string `json:"jenis"       binding:"required,oneof=penambahan pengurangan"`
	Jumlah     int    `json:"jumlah"      binding:"required,min=1"`
	Keterangan string `json:"keterangan"` // opsional
}

// Response
type StokHistoryResponse struct {
	IDHistory   string    `json:"id_history"`
	IDProduk    string    `json:"id_produk"`
	NamaProduk  string    `json:"nama_produk"`
	IDUser      string    `json:"id_user"`
	NamaUser    string    `json:"nama_user"`
	Jenis       string    `json:"jenis"`
	Jumlah      int       `json:"jumlah"`
	StokSebelum int       `json:"stok_sebelum"`
	StokSesudah int       `json:"stok_sesudah"`
	Keterangan  string    `json:"keterangan"`
	NilaiRugi   float64   `json:"nilai_rugi"`
	Tanggal     time.Time `json:"tanggal"`
}