package services

import (
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repositories.NewProductRepository(),
	}
}

func (s *ProductService) GetAll() ([]models.Produk, error) {
	return s.productRepo.FindAll()
}

func (s *ProductService) GetByID(id string) (*models.Produk, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(req *models.CreateProdukRequest) (*models.Produk, error) {
	// Set default status
	status := req.StatusProduk
	if status == "" {
		status = "ready"
	}

	// Parse expired date (opsional)
	var expiredDate *time.Time
	if req.ExpiredDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpiredDate)
		if err == nil {
			expiredDate = &parsed
		}
	}

	var imageURL *string
	if req.ImageURL != "" {
		imageURL = &req.ImageURL
	}

	product := &models.Produk{
		IDProduk:     req.IDProduk,
		NamaProduk:   req.NamaProduk,
		HargaModal:   req.HargaModal,
		HargaJual:    req.HargaJual,
		Stok:         req.Stok,
		IDKategori:   req.IDKategori,
		StatusProduk: status,
		ExpiredDate:  expiredDate,
		ImageURL:     imageURL,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Update(id string, req *models.UpdateProdukRequest) (*models.Produk, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.NamaProduk != "" {
		product.NamaProduk = req.NamaProduk
	}
	if req.HargaModal != 0 {
		product.HargaModal = req.HargaModal
	}
	if req.HargaJual != 0 {
		product.HargaJual = req.HargaJual
	}
	if req.Stok != 0 {
		product.Stok = req.Stok
	}
	if req.IDKategori != "" {
		product.IDKategori = req.IDKategori
	}

	if req.StatusProduk != "" {
		product.StatusProduk = req.StatusProduk
	}
	if req.ExpiredDate != "" {
		parsed, err := time.Parse("2006-01-02", req.ExpiredDate)
		if err == nil {
			product.ExpiredDate = &parsed
		}
	}

	if req.ImageURL != "" {
		product.ImageURL = &req.ImageURL
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Delete(id string) error {
	return s.productRepo.Delete(id)
}
