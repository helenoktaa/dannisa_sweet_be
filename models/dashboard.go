package models

// DashboardResponse - 3 info utama yang ditampilkan di dashboard
// sesuai revisi dosen
type DashboardResponse struct {
	// 1. Transaksi pending pembayaran
	TotalPending      int64              `json:"total_pending"`
	TransaksiPending  []TransaksiPending `json:"transaksi_pending"`

	// 2. Produk mendekati expired (dalam 7 hari)
	TotalMendekatiExpired  int64           `json:"total_mendekati_expired"`
	ProdukMendekatiExpired []ProdukExpired `json:"produk_mendekati_expired"`

	// 3. Produk stok menipis (stok <= 5)
	TotalStokMenipis  int64        `json:"total_stok_menipis"`
	ProdukStokMenipis []ProdukStok `json:"produk_stok_menipis"`
}

// TransaksiPending - transaksi yang belum lunas
type TransaksiPending struct {
	IDTransaksi      string  `json:"id_transaksi"`
	NamaCustomer     string  `json:"nama_customer"`
	JumlahBayar      float64 `json:"jumlah_bayar"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	TanggalTransaksi string  `json:"tanggal_transaksi"`
	HariMenunggu     int     `json:"hari_menunggu"`    // sudah berapa hari pending
	SudahLewat3Hari  bool    `json:"sudah_lewat_3_hari"` // warning jika > 3 hari
}

// ProdukExpired - produk yang expired dalam 7 hari ke depan
type ProdukExpired struct {
	IDProduk    string `json:"id_produk"`
	NamaProduk  string `json:"nama_produk"`
	Stok        int    `json:"stok"`
	ExpiredDate string `json:"expired_date"` // format: 2006-01-02
	SisaHari    int    `json:"sisa_hari"`
}

// ProdukStok - produk yang stoknya menipis
type ProdukStok struct {
	IDProduk     string `json:"id_produk"`
	NamaProduk   string `json:"nama_produk"`
	Stok         int    `json:"stok"`
	StatusProduk string `json:"status_produk"` // ready_stock / pre_order
}

// DashboardHarian - data harian
type DashboardHarian struct {
    TotalPendingLewat3Hari int64              `json:"total_pending_lewat_3_hari"`
    TotalLunasHariIni      int64              `json:"total_lunas_hari_ini"`
    KeuntunganBersih       float64            `json:"keuntungan_bersih"`
    TotalOmzet             float64            `json:"total_omzet"`
    TotalModal             float64            `json:"total_modal"`
    TotalTransaksi         int64              `json:"total_transaksi"`
    TransaksiTerbaru       []TransaksiTerbaru `json:"transaksi_terbaru"`
}

type TransaksiTerbaru struct {
    IDTransaksi      string  `json:"id_transaksi"`
    NamaCustomer     string  `json:"nama_customer"`
    TanggalTransaksi string  `json:"tanggal_transaksi"`
    TotalItem        int     `json:"total_item"`
    JumlahBayar      float64 `json:"jumlah_bayar"`
    MetodePembayaran string  `json:"metode_pembayaran"`
    StatusPembayaran string  `json:"status_pembayaran"`
}