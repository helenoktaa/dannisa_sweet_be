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

// Create membuat transaksi baru beserta detail itemnya
func (s *TransaksiService) Create(req models.CreateTransaksiRequest) (*models.TransaksiResponse, error) {
	// 1. Validasi setiap produk — cek stok cukup
	for _, item := range req.Detail {
		produk, err := s.produkRepo.FindByID(item.IDProduk)
		if err != nil {
			return nil, fmt.Errorf("produk %s tidak ditemukan", item.IDProduk)
		}
		if produk.Stok < item.Qty {
			return nil, fmt.Errorf("stok produk %s tidak cukup (stok: %d, diminta: %d)",
				produk.NamaProduk, produk.Stok, item.Qty)
		}
	}

	// 2. Build detail transaksi
	var details []models.DetailTransaksi
	var totalPenjualan float64

	for _, item := range req.Detail {
		produk, _ := s.produkRepo.FindByID(item.IDProduk)

		detail := models.DetailTransaksi{
			IDTransaksi: req.IDTransaksi,
			IDProduk:    item.IDProduk,
			Qty:         item.Qty,
			HargaJual:   produk.HargaJual, // snapshot harga saat transaksi
		}
		details = append(details, detail)
		totalPenjualan += produk.HargaJual * float64(item.Qty)
	}

	// 3. Build transaksi
	transaksi := &models.Transaksi{
		IDTransaksi:      req.IDTransaksi,
		TanggalTransaksi: time.Now(),
		NamaCustomer:     req.NamaCustomer,
		JumlahBayar:      req.JumlahBayar,
		MetodePembayaran: req.MetodePembayaran,
		StatusPembayaran: "Pending", // confirm by admin/kasir
		IDUser:           req.IDUser,
		Detail:           details,
	}

	// 4. Simpan ke DB (repository handle update stok otomatis)
	if err := s.transaksiRepo.Create(transaksi); err != nil {
		return nil, errors.New("gagal menyimpan transaksi")
	}

	// 5. Ambil data lengkap dengan relasi
	saved, err := s.transaksiRepo.FindByID(transaksi.IDTransaksi)
	if err != nil {
		return nil, err
	}

	return s.buildResponse(saved), nil
}

// GetAll mengambil semua transaksi dengan filter tanggal opsional
func (s *TransaksiService) GetAll(tanggalMulai, tanggalAkhir string) (*models.TransaksiListResponse, error) {
	transaksis, err := s.transaksiRepo.FindAll(tanggalMulai, tanggalAkhir)
	if err != nil {
		return nil, err
	}

	var responses []models.TransaksiResponse
	for _, t := range transaksis {
		responses = append(responses, *s.buildResponse(&t))
	}

	return &models.TransaksiListResponse{
		Data:  responses,
		Total: int64(len(responses)),
	}, nil
}

// GetByID mengambil satu transaksi
func (s *TransaksiService) GetByID(id string) (*models.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}
	return s.buildResponse(transaksi), nil
}

// GetStruk generate data struk dari transaksi
func (s *TransaksiService) GetInvoice(id string) (*models.InvoiceResponse, error) {
    transaksi, err := s.transaksiRepo.FindByID(id)
    if err != nil {
        return nil, errors.New("transaksi tidak ditemukan")
    }

    resp := s.buildResponse(transaksi)

    return &models.InvoiceResponse{
        IDTransaksi:      transaksi.IDTransaksi,
        TanggalTransaksi: transaksi.TanggalTransaksi,
        NamaCustomer:     transaksi.NamaCustomer,
        NamaKasir:        transaksi.User.NamaUser,
        MetodePembayaran: transaksi.MetodePembayaran,
        StatusPembayaran: transaksi.StatusPembayaran,
        Detail:           resp.Detail,
        TotalItem:        resp.TotalItem,
        TotalPenjualan:   resp.TotalPenjualan,
        JumlahBayar:      transaksi.JumlahBayar,
        Kembalian:        transaksi.JumlahBayar - resp.TotalPenjualan,
        InfoPembayaran: models.InfoPembayaran{
            NamaRekening: transaksi.User.NamaUser,
            NoRekening:   transaksi.User.RekPembayaran,
            WhatsApp:     transaksi.User.Whatsapp,
            Catatan:      "Mohon transfer sesuai nominal dan konfirmasi via WhatsApp",
        },
    }, nil
}

// GetLaporan laporan penjualan dengan kalkulasi modal dan laba
func (s *TransaksiService) GetLaporan(req models.LaporanRequest) (*models.LaporanResponse, error) {
	transaksis, err := s.transaksiRepo.FindAll(req.TanggalMulai, req.TanggalAkhir)
	if err != nil {
		return nil, err
	}

	var responses []models.TransaksiResponse
	var totalPenjualan, totalModal float64

	for _, t := range transaksis {
		resp := s.buildResponse(&t)
		responses = append(responses, *resp)
		totalPenjualan += resp.TotalPenjualan

		// Hitung modal dari harga_modal produk
		for _, d := range t.Detail {
			totalModal += d.Produk.HargaModal * float64(d.Qty)
		}
	}

	return &models.LaporanResponse{
		TanggalMulai:   req.TanggalMulai,
		TanggalAkhir:   req.TanggalAkhir,
		TotalTransaksi: int64(len(transaksis)),
		TotalPenjualan: totalPenjualan,
		TotalModal:     totalModal,
		TotalLaba:      totalPenjualan - totalModal,
		Transaksis:     responses,
	}, nil
}

// UpdateStatus update status pembayaran
func (s *TransaksiService) UpdateStatus(id string, req models.UpdateStatusPembayaranRequest) (*models.TransaksiResponse, error) {
	if _, err := s.transaksiRepo.FindByID(id); err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if err := s.transaksiRepo.UpdateStatus(id, req.StatusPembayaran); err != nil {
		return nil, errors.New("gagal update status")
	}

	updated, _ := s.transaksiRepo.FindByID(id)
	return s.buildResponse(updated), nil
}

// buildResponse helper — konversi Transaksi model ke TransaksiResponse
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
		MetodePembayaran: t.MetodePembayaran,
		StatusPembayaran: t.StatusPembayaran,
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