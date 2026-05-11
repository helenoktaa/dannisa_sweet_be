package models

// Kategori menyimpan kategori produk (misal: Bolen, Kue Kering, Cake, dll)
type Kategori struct {
	IDKategori   int    `gorm:"primaryKey;autoIncrement"  json:"id_kategori"`
	NamaKategori string `gorm:"not null;size:50;unique"   json:"nama_kategori"`

	// Relasi: satu kategori punya banyak produk
	Produks []Produk `gorm:"foreignKey:IDKategori" json:"produks,omitempty"`
}

// DTO
type CreateKategoriRequest struct {
	NamaKategori string `json:"nama_kategori" binding:"required"`
}

type UpdateKategoriRequest struct {
	NamaKategori string `json:"nama_kategori" binding:"required"`
}

// Response
type KategoriResponse struct {
	IDKategori   int    `json:"id_kategori"`
	NamaKategori string `json:"nama_kategori"`
}

type KategoriWithProdukResponse struct {
	IDKategori   int               `json:"id_kategori"`
	NamaKategori string            `json:"nama_kategori"`
	Produks      []ProdukResponse  `json:"produks"`
}