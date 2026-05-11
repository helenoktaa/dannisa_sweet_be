package models

// DetailTransaksi menyimpan item-item dalam satu transaksi
// Menggunakan composite primary key: id_transaksi + id_produk
// TIDAK ada kolom total_harga karena bisa dihitung dari qty * harga_jual
type DetailTransaksi struct {
	IDTransaksi string  `gorm:"primaryKey;size:30;index" json:"id_transaksi"`
	IDProduk    string  `gorm:"primaryKey;size:20;index" json:"id_produk"`
	Qty         int     `gorm:"not null;default:1"       json:"qty"`
	HargaJual   float64 `gorm:"not null"                 json:"harga_jual"` // snapshot harga saat transaksi terjadi

	// Relasi
	Transaksi Transaksi `gorm:"foreignKey:IDTransaksi" json:"transaksi,omitempty"`
	Produk    Produk    `gorm:"foreignKey:IDProduk"    json:"produk,omitempty"`
}

// DTO - Create Detail Transaksi
// Dipakai sebagai bagian dari CreateTransaksiRequest (bukan endpoint terpisah)
type CreateDetailTransaksiRequest struct {
	IDProduk  string `json:"id_produk"  binding:"required"`
	Qty       int    `json:"qty"        binding:"required,min=1"`
}

// Response - satu item detail transaksi
// SubTotal dihitung di backend: qty * harga_jual
// TIDAK disimpan di DB sesuai review dosen
type DetailTransaksiResponse struct {
	IDTransaksi string         `json:"id_transaksi"`
	IDProduk    string         `json:"id_produk"`
	Qty         int            `json:"qty"`
	HargaJual   float64        `json:"harga_jual"`
	SubTotal    float64        `json:"sub_total"`    // dihitung: qty * harga_jual (bukan dari DB)
	Produk      ProdukResponse `json:"produk"`
}

// CartItem - keranjang belanja sementara sebelum checkout
// Disimpan di memory/state Flutter, TIDAK perlu tabel di DB
// karena transaksi POS langsung diproses saat itu juga
type CartItem struct {
	IDProduk   string  `json:"id_produk"`
	NamaProduk string  `json:"nama_produk"`
	HargaJual  float64 `json:"harga_jual"`
	Qty        int     `json:"qty"`
	SubTotal   float64 `json:"sub_total"` // qty * harga_jual
}

// CartRequest & CartResponse - dipakai di Flutter (state management)
// BUKAN endpoint API, hanya untuk perhitungan lokal di Flutter
type CartResponse struct {
	Items          []CartItem `json:"items"`
	TotalItem      int        `json:"total_item"`
	TotalPenjualan float64    `json:"total_penjualan"` // SUM semua sub_total
}