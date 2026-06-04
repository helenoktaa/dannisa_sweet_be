package services

import (
	"fmt"
	"time"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type ProductService struct {
	productRepo     *repositories.ProductRepository
	markdownService *MarkdownPricingService
	historyService  *StokHistoryService
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo:     repositories.NewProductRepository(),
		markdownService: NewMarkdownPricingService(),
		historyService:  NewStokHistoryService(),
	}
}

func (s *ProductService) GetAll(statusProduk string) ([]models.ProdukResponse, error) {
	products, err := s.productRepo.FindAll(statusProduk)
	if err != nil {
		return nil, err
	}

	var responses []models.ProdukResponse
	for _, p := range products {
		responses = append(responses, s.buildProdukResponse(p))
	}

	return responses, nil
}

func (s *ProductService) GetByID(id string) (*models.ProdukResponse, error) {
	p, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := s.buildProdukResponse(*p)
	return &resp, nil
}

func (s *ProductService) Create(req *models.CreateProdukRequest, idUser string) (*models.Produk, error) {
	status := req.StatusProduk
	if status == "" {
		status = "ready"
	}

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

	// Catat stok awal ke history jika stok > 0
	if req.Stok > 0 {
		err2 := s.historyService.CatatSaja(
			product.IDProduk,
			idUser,
			req.Stok,
			"Stok awal produk baru",
		)
		if err2 != nil {
			fmt.Println("=== ERROR CATAT HISTORY:", err2, "===")
		}
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

func (s *ProductService) buildProdukResponse(p models.Produk) models.ProdukResponse {
	resp := models.ProdukResponse{
		IDProduk:     p.IDProduk,
		NamaProduk:   p.NamaProduk,
		HargaModal:   p.HargaModal,
		HargaJual:    p.HargaJual,
		Stok:         p.Stok,
		StatusProduk: p.StatusProduk,
		ExpiredDate:  p.ExpiredDate,
		ImageURL:     p.ImageURL,
		Kategori: models.KategoriResponse{
			IDKategori:   p.Kategori.IDKategori,
			NamaKategori: p.Kategori.NamaKategori,
		},
	}

	harga := s.markdownService.HitungHargaEfektif(p)
	if harga.SumberDiskon != "tidak_ada" {
		resp.HargaDiskon = &harga.HargaDiskon
		resp.PorsenDiskon = &harga.PorsenDiskon
		resp.SumberDiskon = &harga.SumberDiskon
	}

	return resp
}
