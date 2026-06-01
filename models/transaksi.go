package models

import "time"

const (
	JenisReadyStock = "ready_stock"
	JenisPreOrder   = "pre_order"

	StatusMenungguDiproses = "menunggu_diproses"
	StatusSedangDibuat     = "sedang_dibuat"
	StatusSedangDiantar    = "sedang_diantar"
	StatusPesananDiterima  = "pesanan_diterima"
	StatusSelesai          = "selesai"
	StatusDibatalkan       = "dibatalkan"
)

// Transaksi menyimpan data transaksi penjualan Dannisa Sweet
type Transaksi struct {
	IDTransaksi      string    `gorm:"primaryKey;size:30"                        json:"id_transaksi"`
	TanggalTransaksi time.Time `gorm:"not null"                                  json:"tanggal_transaksi"`
	NamaCustomer     string    `gorm:"not null;size:100"                         json:"nama_customer"`
	JumlahBayar      float64   `gorm:"not null"                                  json:"jumlah_bayar"`
	MetodePembayaran string    `gorm:"not null;size:30"                          json:"metode_pembayaran"` // Tunai / Transfer / QRIS
	StatusPembayaran string    `gorm:"not null;size:20;default:'Pending'"        json:"status_pembayaran"` // Pending / Lunas
	IDUser           string    `gorm:"not null;size:20;index"                    json:"id_user"`
	JenisOrder       string    `gorm:"not null;size:20;default:'ready_stock'"   json:"jenis_order"`
	StatusOrder      string    `gorm:"not null;size:30;default:'selesai'"        json:"status_order"`
	Catatan          string    `gorm:"size:500"                                  json:"catatan"`

	// Relasi
	User   User              `gorm:"foreignKey:IDUser;references:IDUser"    json:"user,omitempty"`
	Detail []DetailTransaksi `gorm:"foreignKey:IDTransaksi" json:"detail,omitempty"`
}

// DTO - Create Transaksi
// Detail langsung disertakan saat buat transaksi (satu request)
type CreateTransaksiRequest struct {
	IDTransaksi      string `json:"id_transaksi"`
	NamaCustomer     string `json:"nama_customer"     binding:"required"`
	MetodePembayaran string `json:"metode_pembayaran" binding:"required,oneof=Tunai Transfer QRIS"`
	// JumlahBayar tidak wajib saat create — diisi saat update status
	JumlahBayar float64                        `json:"jumlah_bayar"`
	IDUser      string                         `json:"id_user"`
	Detail      []CreateDetailTransaksiRequest `json:"detail"            binding:"required,min=1,dive"`
	JenisOrder  string                         `json:"jenis_order" binding:"omitempty,oneof=ready_stock pre_order"`
	Catatan     string                         `json:"catatan"`
}

// DTO - Update Status Pembayaran
type UpdateStatusPembayaranRequest struct {
	StatusPembayaran string  `json:"status_pembayaran" binding:"required,oneof=Pending Lunas"`
	JumlahBayar      float64 `json:"jumlah_bayar"      binding:"omitempty,min=0"`
}

// DTO - Update Status Order (khusus pre order)
type UpdateStatusOrderRequest struct {
	StatusOrder string `json:"status_order" binding:"required,oneof=sedang_dibuat sedang_diantar pesanan_diterima selesai dibatalkan"`
	Catatan     string `json:"catatan"`
}

// Response - satu transaksi
// TotalPenjualan dihitung di backend: SUM(qty * harga_jual)
// TIDAK disimpan di DB karena bisa dihitung dari detail
type TransaksiResponse struct {
	IDTransaksi      string                    `json:"id_transaksi"`
	TanggalTransaksi time.Time                 `json:"tanggal_transaksi"`
	NamaCustomer     string                    `json:"nama_customer"`
	JumlahBayar      float64                   `json:"jumlah_bayar"`
	MetodePembayaran string                    `json:"metode_pembayaran"`
	StatusPembayaran string                    `json:"status_pembayaran"`
	JenisOrder       string                    `json:"jenis_order"`
	StatusOrder      string                    `json:"status_order"`
	Catatan          string                    `json:"catatan"`
	User             UserResponse              `json:"user"`
	Detail           []DetailTransaksiResponse `json:"detail"`

	
	TotalItem      int     `json:"total_item"`
	TotalPenjualan float64 `json:"total_penjualan"` // SUM(qty * harga_jual)
}

// Response - list transaksi
type TransaksiListResponse struct {
	Data  []TransaksiResponse `json:"data"`
	Total int64               `json:"total"`
}

// DTO & Response Laporan
// Laporan TIDAK punya tabel sendiri — dihitung dari JOIN query
type LaporanRequest struct {
	TanggalMulai string `form:"tanggal_mulai" binding:"required"` // format: 2024-01-01
	TanggalAkhir string `form:"tanggal_akhir" binding:"required"` // format: 2024-12-31
}

type LaporanResponse struct {
	TanggalMulai   string              `json:"tanggal_mulai"`
	TanggalAkhir   string              `json:"tanggal_akhir"`
	TotalTransaksi int64               `json:"total_transaksi"`
	TotalPenjualan float64             `json:"total_penjualan"` // SUM(qty * harga_jual)
	TotalModal     float64             `json:"total_modal"`     // SUM(qty * harga_modal dari tabel produk)
	TotalLaba      float64             `json:"total_laba"`      // total_penjualan - total_modal
	Transaksis     []TransaksiResponse `json:"transaksis"`
}

type InvoiceResponse struct {
	IDTransaksi      string                    `json:"id_transaksi"`
	TanggalTransaksi time.Time                 `json:"tanggal_transaksi"`
	NamaCustomer     string                    `json:"nama_customer"`
	NamaKasir        string                    `json:"nama_kasir"`
	MetodePembayaran string                    `json:"metode_pembayaran"`
	StatusPembayaran string                    `json:"status_pembayaran"`
	Detail           []DetailTransaksiResponse `json:"detail"`
	TotalItem        int                       `json:"total_item"`
	TotalPenjualan   float64                   `json:"total_penjualan"`
	JumlahBayar      float64                   `json:"jumlah_bayar"`
	Kembalian        float64                   `json:"kembalian"`
	InfoPembayaran   InfoPembayaran            `json:"info_pembayaran"`
}

type InfoPembayaran struct {
	NamaRekening string `json:"nama_rekening"` // Anisa Dian Utami
	NoRekening   string `json:"no_rekening"`   // BCA 8880587898
	WhatsApp     string `json:"whatsapp"`      // 085156194878
	Catatan      string `json:"catatan"`       // "Mohon transfer sesuai nominal dan konfirmasi via WA"
}
