package services

import (
	"errors"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *UserService) GetAll() ([]models.UserResponse, error) {
	var users []models.User
	if err := repositories.NewUserRepository().FindAll(&users); err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, u := range users {
		responses = append(responses, models.UserResponse{
			IDUser:        u.IDUser,
			NamaUser:      u.NamaUser,
			Email:         u.Email,
			RekPembayaran: u.RekPembayaran,
			Whatsapp:      u.Whatsapp,
			Jabatan: models.JabatanResponse{
				IDJabatan:   u.Jabatan.IDJabatan,
				NamaJabatan: u.Jabatan.NamaJabatan,
				Gaji:        u.Jabatan.Gaji,
			},
		})
	}
	return responses, nil
}

func (s *UserService) GetByID(id string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &models.UserResponse{
		IDUser:        user.IDUser,
		NamaUser:      user.NamaUser,
		Email:         user.Email,
		RekPembayaran: user.RekPembayaran,
		Whatsapp:      user.Whatsapp,
		Jabatan: models.JabatanResponse{
			IDJabatan:   user.Jabatan.IDJabatan,
			NamaJabatan: user.Jabatan.NamaJabatan,
			Gaji:        user.Jabatan.Gaji,
		},
	}, nil
}

func (s *UserService) Update(id string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.NamaUser != "" {
		user.NamaUser = req.NamaUser
	}
	if req.RekPembayaran != "" {
		user.RekPembayaran = req.RekPembayaran
	}
	if req.Whatsapp != "" {
		user.Whatsapp = req.Whatsapp
	}
	if req.IDJabatan != "" {
		user.IDJabatan = req.IDJabatan
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("gagal update user")
	}

	updated, _ := s.userRepo.FindByID(id)
	return &models.UserResponse{
		IDUser:        updated.IDUser,
		NamaUser:      updated.NamaUser,
		Email:         updated.Email,
		RekPembayaran: updated.RekPembayaran,
		Whatsapp:      updated.Whatsapp,
		Jabatan: models.JabatanResponse{
			IDJabatan:   updated.Jabatan.IDJabatan,
			NamaJabatan: updated.Jabatan.NamaJabatan,
			Gaji:        updated.Jabatan.Gaji,
		},
	}, nil
}

func (s *UserService) Delete(id string) error {
	if _, err := s.userRepo.FindByID(id); err != nil {
		return errors.New("user tidak ditemukan")
	}
	return s.userRepo.Delete(id)
}