package services

import (
	"errors"
	"fmt"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type KategoriService struct {
	kategoriRepo *repositories.KategoriRepository
}

func NewKategoriService() *KategoriService {
	return &KategoriService{
		kategoriRepo: repositories.NewKategoriRepository(),
	}
}

func (s *KategoriService) GetAll() ([]models.KategoriResponse, error) {
	kategoris, err := s.kategoriRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.KategoriResponse
	for _, k := range kategoris {
		responses = append(responses, models.KategoriResponse{
			IDKategori:   k.IDKategori,
			NamaKategori: k.NamaKategori,
		})
	}
	return responses, nil
}

func (s *KategoriService) GetByID(id string) (*models.KategoriResponse, error) {
	kategori, err := s.kategoriRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	return &models.KategoriResponse{
		IDKategori:   kategori.IDKategori,
		NamaKategori: kategori.NamaKategori,
	}, nil
}

func (s *KategoriService) Create(req models.CreateKategoriRequest) (*models.KategoriResponse, error) {
	// Auto generate ID
	kategoris, _ := s.kategoriRepo.FindAll()
	id := fmt.Sprintf("KDS%03d", len(kategoris)+1)

	kategori := &models.Kategori{
		IDKategori:   id,
		NamaKategori: req.NamaKategori,
	}

	if err := s.kategoriRepo.Create(kategori); err != nil {
		return nil, errors.New("gagal membuat kategori")
	}

	return &models.KategoriResponse{
		IDKategori:   kategori.IDKategori,
		NamaKategori: kategori.NamaKategori,
	}, nil
}

func (s *KategoriService) Update(id string, req models.UpdateKategoriRequest) (*models.KategoriResponse, error) {
	kategori, err := s.kategoriRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	if req.NamaKategori != "" {
		kategori.NamaKategori = req.NamaKategori
	}

	if err := s.kategoriRepo.Update(kategori); err != nil {
		return nil, errors.New("gagal update kategori")
	}

	return &models.KategoriResponse{
		IDKategori:   kategori.IDKategori,
		NamaKategori: kategori.NamaKategori,
	}, nil
}

func (s *KategoriService) Delete(id string) error {
	if _, err := s.kategoriRepo.FindByID(id); err != nil {
		return errors.New("kategori tidak ditemukan")
	}
	return s.kategoriRepo.Delete(id)
}