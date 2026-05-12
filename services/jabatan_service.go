package services

import (
	"errors"
	"fmt"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type JabatanService struct {
	jabatanRepo *repositories.JabatanRepository
}

func NewJabatanService() *JabatanService {
	return &JabatanService{
		jabatanRepo: repositories.NewJabatanRepository(),
	}
}

func (s *JabatanService) GetAll() ([]models.JabatanResponse, error) {
	jabatans, err := s.jabatanRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.JabatanResponse
	for _, j := range jabatans {
		responses = append(responses, models.JabatanResponse{
			IDJabatan:   j.IDJabatan,
			NamaJabatan: j.NamaJabatan,
			Gaji:        j.Gaji,
		})
	}
	return responses, nil
}

func (s *JabatanService) GetByID(id string) (*models.JabatanResponse, error) {
	jabatan, err := s.jabatanRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("jabatan tidak ditemukan")
	}

	return &models.JabatanResponse{
		IDJabatan:   jabatan.IDJabatan,
		NamaJabatan: jabatan.NamaJabatan,
		Gaji:        jabatan.Gaji,
	}, nil
}

func (s *JabatanService) Create(req models.CreateJabatanRequest) (*models.JabatanResponse, error) {
	// Auto generate ID
	jabatans, _ := s.jabatanRepo.FindAll()
	id := fmt.Sprintf("JAB%03d", len(jabatans)+1)

	jabatan := &models.Jabatan{
		IDJabatan:   id,
		NamaJabatan: req.NamaJabatan,
		Gaji:        req.Gaji,
	}

	if err := s.jabatanRepo.Create(jabatan); err != nil {
		return nil, errors.New("gagal membuat jabatan")
	}

	return &models.JabatanResponse{
		IDJabatan:   jabatan.IDJabatan,
		NamaJabatan: jabatan.NamaJabatan,
		Gaji:        jabatan.Gaji,
	}, nil
}

func (s *JabatanService) Update(id string, req models.UpdateJabatanRequest) (*models.JabatanResponse, error) {
	jabatan, err := s.jabatanRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("jabatan tidak ditemukan")
	}

	if req.NamaJabatan != "" {
		jabatan.NamaJabatan = req.NamaJabatan
	}
	if req.Gaji > 0 {
		jabatan.Gaji = req.Gaji
	}

	if err := s.jabatanRepo.Update(jabatan); err != nil {
		return nil, errors.New("gagal update jabatan")
	}

	return &models.JabatanResponse{
		IDJabatan:   jabatan.IDJabatan,
		NamaJabatan: jabatan.NamaJabatan,
		Gaji:        jabatan.Gaji,
	}, nil
}

func (s *JabatanService) Delete(id string) error {
	if _, err := s.jabatanRepo.FindByID(id); err != nil {
		return errors.New("jabatan tidak ditemukan")
	}
	return s.jabatanRepo.Delete(id)
}