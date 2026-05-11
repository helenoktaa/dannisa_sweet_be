package services

import (
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

// Get all produk
func (s *ProductService) GetAll() ([]models.Produk, error) {
	return s.productRepo.FindAll()
}

// Get produk by ID
func (s *ProductService) GetByID(id string) (*models.Produk, error) {
	return s.productRepo.FindByID(id)
}

// Create produk
func (s *ProductService) Create(req *models.CreateProdukRequest) (*models.Produk, error) {
	product := &models.Produk{
		IDProduk:    req.IDProduk,
		NamaProduk:  req.NamaProduk,
		HargaModal:  req.HargaModal,
		HargaJual:   req.HargaJual,
		Stok:        req.Stok,
		IDKategori:  req.IDKategori,
	}

	err := s.productRepo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Update produk
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

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Delete produk
func (s *ProductService) Delete(id string) error {
	return s.productRepo.Delete(id)
}