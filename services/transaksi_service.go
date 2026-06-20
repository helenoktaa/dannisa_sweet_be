package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type TransaksiService struct {
	transaksiRepo *repositories.TransaksiRepository
	produkRepo    *repositories.ProductRepository
}

func NewTransaksiService() *TransaksiService {
	return &TransaksiService{
		transaksiRepo: repositories.NewTransaksiRepository(),
		produkRepo:    repositories.NewProductRepository(),
	}
}

// ─────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────
func (s *TransaksiService) Create(req models.CreateTransaksiRequest) (*models.TransaksiResponse, error) {
	lastNumber, err := s.transaksiRepo.GetLastNumber()
	if err != nil {
		return nil, errors.New("gagal generate ID transaksi")
	}
	req.IDTransaksi = fmt.Sprintf("TDS%04d", lastNumber+1)

	jenisOrder := models.JenisReadyStock
	statusOrder := models.StatusSelesai

	for _, item := range req.Detail {
		produk, err := s.produkRepo.FindByID(item.IDProduk)
		if err != nil {
			return nil, fmt.Errorf("produk %s tidak ditemukan", item.IDProduk)
		}

		if produk.StatusProduk == "preorder" {
			jenisOrder = models.JenisPreOrder
			statusOrder = models.StatusMenungguDiproses
		} else {
			// Hanya cek stok untuk produk ready stock
			if produk.StatusProduk != "preorder" && produk.Stok < item.Qty {
				return nil, fmt.Errorf("stok produk %s tidak cukup (stok: %d, diminta: %d)",
					produk.NamaProduk, produk.Stok, item.Qty)
			}
		}
	}

	var details []models.DetailTransaksi
	var totalHarga float64

	for _, item := range req.Detail {
		produk, _ := s.produkRepo.FindByID(item.IDProduk)
		totalHarga += produk.HargaJual * float64(item.Qty)
		details = append(details, models.DetailTransaksi{
			IDTransaksi: req.IDTransaksi,
			IDProduk:    item.IDProduk,
			Qty:         item.Qty,
			HargaJual:   produk.HargaJual,
		})
	}

	statusPembayaran := "Pending"
	var tanggalLunas *time.Time // ← pointer, nil kalau belum lunas

	if jenisOrder == models.JenisReadyStock &&
		req.MetodePembayaran == "Tunai" &&
		req.JumlahBayar >= totalHarga {
		statusPembayaran = "Lunas"
		now := time.Now()
		tanggalLunas = &now // ← set tanggal lunas sekarang
	}

	transaksi := &models.Transaksi{
		IDTransaksi:      req.IDTransaksi,
		TanggalTransaksi: time.Now(),
		NamaCustomer:     req.NamaCustomer,
		JumlahBayar:      req.JumlahBayar,
		MetodePembayaran: req.MetodePembayaran,
		StatusPembayaran: statusPembayaran,
		IDUser:           req.IDUser,
		JenisOrder:       jenisOrder,
		StatusOrder:      statusOrder,
		Catatan:          req.Catatan,
		Detail:           details,
		TanggalLunas:     tanggalLunas,
	}
	if err := s.transaksiRepo.Create(transaksi); err != nil {
		return nil, errors.New("gagal menyimpan transaksi")
	}

	saved, err := s.transaksiRepo.FindByID(transaksi.IDTransaksi)
	if err != nil {
		return nil, err
	}
	return s.buildResponse(saved), nil
}

// ─────────────────────────────────────────────
// GetAll
// ─────────────────────────────────────────────
func (s *TransaksiService) GetAll(tanggalMulai, tanggalAkhir, status string) ([]models.TransaksiResponse, error) {
	transaksis, err := s.transaksiRepo.FindAll(tanggalMulai, tanggalAkhir, status)
	if err != nil {
		return nil, err
	}

	var result []models.TransaksiResponse
	for _, t := range transaksis {
		result = append(result, *s.buildResponse(&t))
	}
	return result, nil
}

// ─────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────
func (s *TransaksiService) GetByID(id string) (*models.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}
	return s.buildResponse(transaksi), nil
}

// ─────────────────────────────────────────────
// GetInvoice
// ─────────────────────────────────────────────
func (s *TransaksiService) GetInvoice(id string) (*models.InvoiceResponse, error) {
	transaksi, err := s.transaksiRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	resp := s.buildResponse(transaksi)

	totalDibayar := transaksi.JumlahBayar + transaksi.JumlahDp
	kembalian := totalDibayar - resp.TotalPenjualan
	if kembalian < 0 {
		kembalian = 0
	}

	return &models.InvoiceResponse{
		IDTransaksi:      transaksi.IDTransaksi,
		TanggalTransaksi: transaksi.TanggalTransaksi,
		NamaCustomer:     transaksi.NamaCustomer,
		NamaKasir:        transaksi.User.NamaUser,
		MetodePembayaran: transaksi.MetodePembayaran,
		StatusPembayaran: transaksi.StatusPembayaran,
		TanggalLunas:     transaksi.TanggalLunas,
		JenisOrder:       transaksi.JenisOrder,
		JumlahDp:         transaksi.JumlahDp,
		Detail:           resp.Detail,
		TotalItem:        resp.TotalItem,
		TotalPenjualan:   resp.TotalPenjualan,
		JumlahBayar:      transaksi.JumlahBayar,
		Kembalian:        transaksi.JumlahBayar - resp.TotalPenjualan,
		InfoPembayaran: models.InfoPembayaran{
			NamaRekening: "Anisa Dian Utami",
			NoRekening:   "BCA 8880587898",
			WhatsApp:     "085156194878",
			Catatan:      "Mohon transfer sesuai nominal dan konfirmasi via WA",
		},
	}, nil
}

// ─────────────────────────────────────────────
// GetLaporan
// ─────────────────────────────────────────────
func (s *TransaksiService) GetLaporan(req models.LaporanRequest) (*models.LaporanResponse, error) {
	transaksis, err := s.transaksiRepo.FindAll(req.TanggalMulai, req.TanggalAkhir, "")
	if err != nil {
		return nil, err
	}

	var totalModal, totalPenjualan, totalLaba float64
	var totalTransaksi int64
	var transaksiResponses []models.TransaksiResponse

	for _, t := range transaksis {
		if t.StatusPembayaran != "Lunas" {
			continue
		}
		totalTransaksi++
		resp := s.buildResponse(&t)
		for _, d := range t.Detail {
			totalModal += d.Produk.HargaModal * float64(d.Qty)
			totalPenjualan += d.HargaJual * float64(d.Qty)
		}
		transaksiResponses = append(transaksiResponses, *resp)
	}
	totalLaba = totalPenjualan - totalModal

	return &models.LaporanResponse{
		TanggalMulai:   req.TanggalMulai,
		TanggalAkhir:   req.TanggalAkhir,
		TotalTransaksi: totalTransaksi,
		TotalModal:     totalModal,
		TotalPenjualan: totalPenjualan,
		TotalLaba:      totalLaba,
		Transaksis:     transaksiResponses,
	}, nil
}

// ─────────────────────────────────────────────
// UpdateStatus (status pembayaran)
// Stok sudah dikurangi saat Create, tidak dikurangi lagi di sini
// ─────────────────────────────────────────────
func (s *TransaksiService) UpdateStatus(id string, req models.UpdateStatusPembayaranRequest) (*models.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaksi.StatusPembayaran == "Lunas" {
		return nil, errors.New("transaksi sudah lunas")
	}

	total := s.hitungTotal(transaksi)

	if req.StatusPembayaran == "DP" {
		dp50 := total * 0.5
		if req.JumlahDp < dp50 {
			return nil, fmt.Errorf("DP minimal 50%% dari total (Rp %.0f)", dp50)
		}
		if req.JumlahDp > total {
			return nil, errors.New("DP tidak boleh melebihi total pembayaran")
		}
	}

	jumlahBayarFinal := req.JumlahBayar
	jumlahDpFinal := req.JumlahDp

	if req.StatusPembayaran == "Lunas" {
		sudahBayar := transaksi.JumlahBayar + transaksi.JumlahDp
		sisa := total - sudahBayar
		if req.JumlahBayar < sisa {
			return nil, fmt.Errorf("jumlah bayar kurang, sisa yang harus dibayar: Rp %.0f", sisa)
		}
		jumlahBayarFinal = req.JumlahBayar
		jumlahDpFinal = transaksi.JumlahDp
	}

	if err := s.transaksiRepo.UpdateStatusPembayaran(id, req.StatusPembayaran, jumlahBayarFinal, jumlahDpFinal); err != nil {
		return nil, errors.New("gagal mengupdate status pembayaran")
	}

	return s.GetByID(id)
}

// ─────────────────────────────────────────────
// UpdateStatusOrder (khusus pre order)
// ─────────────────────────────────────────────
func (s *TransaksiService) UpdateStatusOrder(id string, req models.UpdateStatusOrderRequest) (*models.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaksi.JenisOrder != models.JenisPreOrder {
		return nil, errors.New("hanya transaksi pre order yang bisa diupdate status order-nya")
	}

	urutanValid := map[string]string{
		models.StatusMenungguDiproses: models.StatusSedangDibuat,
		models.StatusSedangDibuat:     models.StatusSedangDiantar,
		models.StatusSedangDiantar:    models.StatusPesananDiterima,
		models.StatusPesananDiterima:  models.StatusSelesai,
	}

	if req.StatusOrder != models.StatusDibatalkan {
		nextValid, ok := urutanValid[transaksi.StatusOrder]
		if !ok || req.StatusOrder != nextValid {
			return nil, fmt.Errorf("status tidak valid: dari '%s' hanya bisa ke '%s'",
				transaksi.StatusOrder, urutanValid[transaksi.StatusOrder])
		}
	}

	if err := s.transaksiRepo.UpdateStatusOrder(id, req.StatusOrder, req.Catatan); err != nil {
		return nil, errors.New("gagal mengupdate status order")
	}

	return s.GetByID(id)
}

// ─────────────────────────────────────────────
// GetPreOrderAktif
// ─────────────────────────────────────────────
func (s *TransaksiService) GetPreOrderAktif() ([]models.TransaksiResponse, error) {
	transaksis, err := s.transaksiRepo.FindPreOrderAktif()
	if err != nil {
		return nil, err
	}

	var result []models.TransaksiResponse
	for _, t := range transaksis {
		result = append(result, *s.buildResponse(&t))
	}
	return result, nil
}

// ─────────────────────────────────────────────
// buildResponse (internal helper)
// ─────────────────────────────────────────────
func (s *TransaksiService) buildResponse(t *models.Transaksi) *models.TransaksiResponse {
	var details []models.DetailTransaksiResponse
	var totalItem int
	var totalPenjualan float64

	for _, d := range t.Detail {
		subTotal := d.HargaJual * float64(d.Qty)
		totalItem += d.Qty
		totalPenjualan += subTotal

		details = append(details, models.DetailTransaksiResponse{
			IDTransaksi: d.IDTransaksi,
			IDProduk:    d.IDProduk,
			Qty:         d.Qty,
			HargaJual:   d.HargaJual,
			SubTotal:    subTotal,
			Produk: models.ProdukResponse{
				IDProduk:   d.Produk.IDProduk,
				NamaProduk: d.Produk.NamaProduk,
				HargaModal: d.Produk.HargaModal,
				HargaJual:  d.Produk.HargaJual,
				Stok:       d.Produk.Stok,
				StatusProduk: d.Produk.StatusProduk,
				Kategori: models.KategoriResponse{
					IDKategori:   d.Produk.Kategori.IDKategori,
					NamaKategori: d.Produk.Kategori.NamaKategori,
				},
			},
		})
	}

	return &models.TransaksiResponse{
		IDTransaksi:      t.IDTransaksi,
		TanggalTransaksi: t.TanggalTransaksi,
		NamaCustomer:     t.NamaCustomer,
		JumlahBayar:      t.JumlahBayar,
		JumlahDp:         t.JumlahDp,
		MetodePembayaran: t.MetodePembayaran,
		StatusPembayaran: t.StatusPembayaran,
		JenisOrder:       t.JenisOrder,
		StatusOrder:      t.StatusOrder,
		Catatan:          t.Catatan,
		TanggalLunas:     t.TanggalLunas,
		User: models.UserResponse{
			IDUser:   t.User.IDUser,
			NamaUser: t.User.NamaUser,
			Email:    t.User.Email,
			Jabatan: models.JabatanResponse{
				IDJabatan:   t.User.Jabatan.IDJabatan,
				NamaJabatan: t.User.Jabatan.NamaJabatan,
				Gaji:        t.User.Jabatan.Gaji,
			},
		},
		Detail:         details,
		TotalItem:      totalItem,
		TotalPenjualan: totalPenjualan,
	}
}

// ─────────────────────────────────────────────
// hitungTotal (internal helper)
// ─────────────────────────────────────────────
func (s *TransaksiService) hitungTotal(t *models.Transaksi) float64 {
	var total float64
	for _, d := range t.Detail {
		total += d.HargaJual * float64(d.Qty)
	}
	return total
}
