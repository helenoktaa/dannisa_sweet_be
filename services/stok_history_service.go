package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type StokHistoryService struct {
	historyRepo *repositories.StokHistoryRepository
	produkRepo  *repositories.ProductRepository
	userRepo    *repositories.UserRepository
}

func NewStokHistoryService() *StokHistoryService {
	return &StokHistoryService{
		historyRepo: repositories.NewStokHistoryRepository(),
		produkRepo:  repositories.NewProductRepository(),
		userRepo:    repositories.NewUserRepository(),
	}
}

// Create — catat perubahan stok
func (s *StokHistoryService) Create(
	req models.CreateStokHistoryRequest,
	idUser string,
) (*models.StokHistoryResponse, error) {

	// Cek produk ada
	produk, err := s.produkRepo.FindByID(req.IDProduk)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	stokSebelum := produk.Stok
	var stokSesudah int

	if req.Jenis == "penambahan" {
		stokSesudah = stokSebelum + req.Jumlah
	} else {
		// pengurangan
		if req.Jumlah > stokSebelum {
			return nil, errors.New("jumlah pengurangan melebihi stok yang tersedia")
		}
		stokSesudah = stokSebelum - req.Jumlah
	}

	// Generate ID
	lastNumber, _ := s.historyRepo.GetLastNumber()
	idHistory := fmt.Sprintf("STK%04d", lastNumber+1)

	history := &models.StokHistory{
		IDHistory:   idHistory,
		IDProduk:    req.IDProduk,
		IDUser:      idUser,
		Jenis:       req.Jenis,
		Jumlah:      req.Jumlah,
		StokSebelum: stokSebelum,
		StokSesudah: stokSesudah,
		Keterangan:  req.Keterangan,
		Tanggal:     time.Now(),
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, errors.New("gagal menyimpan history stok")
	}

	// Update stok produk
	produk.Stok = stokSesudah
	if err := s.produkRepo.Update(produk); err != nil {
		return nil, errors.New("gagal update stok produk")
	}

	namaUser := ""
user, err := s.userRepo.FindByID(idUser)
if err == nil {
	namaUser = user.NamaUser
}

	return &models.StokHistoryResponse{
		IDHistory:   history.IDHistory,
		IDProduk:    produk.IDProduk,
		NamaProduk:  produk.NamaProduk,
		IDUser:      idUser,
		NamaUser:    namaUser,
		Jenis:       req.Jenis,
		Jumlah:      req.Jumlah,
		StokSebelum: stokSebelum,
		StokSesudah: stokSesudah,
		Keterangan:  req.Keterangan,
		Tanggal:     history.Tanggal,
	}, nil
}


// GetAll — ambil semua history
func (s *StokHistoryService) GetAll(idProduk, jenis string) ([]models.StokHistoryResponse, error) {
	histories, err := s.historyRepo.FindAll(idProduk, jenis)
	if err != nil {
		return nil, err
	}

	var responses []models.StokHistoryResponse
	for _, h := range histories {
		responses = append(responses, models.StokHistoryResponse{
			IDHistory:   h.IDHistory,
			IDProduk:    h.IDProduk,
			NamaProduk:  h.Produk.NamaProduk,
			IDUser:      h.IDUser,
			NamaUser:    h.User.NamaUser,
			Jenis:       h.Jenis,
			Jumlah:      h.Jumlah,
			StokSebelum: h.StokSebelum,
			StokSesudah: h.StokSesudah,
			Keterangan:  h.Keterangan,
			Tanggal:     h.Tanggal,
		})
	}

	return responses, nil
}