package services

import (
	"errors"
	"fmt"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo        *repositories.UserRepository
	userMenuService *UserMenuService // ← inject service, bukan repo langsung
}

func NewUserService(
	userRepo *repositories.UserRepository,
	userMenuService *UserMenuService,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		userMenuService: userMenuService,
	}
}

// ── Mapping Helper ─────────────────────────────────────────

func toUserResponse(u *models.User) models.UserResponse {
	return models.UserResponse{
		IDUser:        u.IDUser,
		NamaUser:      u.NamaUser,
		Email:         u.Email,
		RekPembayaran: u.RekPembayaran,
		Whatsapp:      u.Whatsapp,
		Jabatan: models.JabatanResponse{
			IDJabatan:   u.Jabatan.IDJabatan,
			NamaJabatan: u.Jabatan.NamaJabatan,
		},
		MenuKeys: models.ExtractMenuKeys(u.UserMenus),
	}
}

// ── Generate User ID ───────────────────────────────────────

func (s *UserService) generateUserID() (string, error) {
	last, err := s.userRepo.GetLastNumber()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("UDS%02d", last+1), nil
}

// ── Get All ────────────────────────────────────────────────

func (s *UserService) GetAll() ([]models.UserResponse, error) {
	var users []models.User
	if err := s.userRepo.FindAll(&users); err != nil {
		return nil, err
	}

	responses := make([]models.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, toUserResponse(&u))
	}
	return responses, nil
}

// ── Get By ID ──────────────────────────────────────────────

func (s *UserService) GetByID(idUser string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	resp := toUserResponse(user)
	return &resp, nil
}

// ── Create ─────────────────────────────────────────────────

func (s *UserService) Create(req models.RegisterRequest) error {
	// Cek email duplikat
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return errors.New("email sudah terdaftar")
	}

	// Generate ID
	id, err := s.generateUserID()
	if err != nil {
		return errors.New("gagal generate ID user")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	user := &models.User{
		IDUser:        id,
		NamaUser:      req.NamaUser,
		Email:         req.Email,
		Password:      string(hashed),
		RekPembayaran: req.RekPembayaran,
		Whatsapp:      req.Whatsapp,
		IDJabatan:     req.IDJabatan,
	}

	if err := s.userRepo.Create(user); err != nil {
		return errors.New("gagal menyimpan user")
	}

	// Delegasi ke UserMenuService
	if len(req.MenuKeys) > 0 {
		if err := s.userMenuService.Replace(user.IDUser, req.MenuKeys); err != nil {
			return err
		}
	}

	return nil
}

// ── Update ─────────────────────────────────────────────────

func (s *UserService) Update(idUser string, req models.UpdateUserRequest) error {
	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return errors.New("user tidak ditemukan")
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
		return errors.New("gagal mengupdate user")
	}

	// Delegasi ke UserMenuService
	if req.MenuKeys != nil {
		if err := s.userMenuService.Replace(idUser, req.MenuKeys); err != nil {
			return err
		}
	}

	return nil
}

// ── Delete ─────────────────────────────────────────────────

func (s *UserService) Delete(idUser string) error {
	_, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Delegasi ke UserMenuService
	if err := s.userMenuService.DeleteByUserID(idUser); err != nil {
		return err
	}

	if err := s.userRepo.Delete(idUser); err != nil {
		return errors.New("gagal menghapus user")
	}

	return nil
}

// ── Update Password ────────────────────────────────────────

func (s *UserService) UpdatePassword(idUser string, req models.UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.PasswordLama),
	); err != nil {
		return errors.New("password lama tidak sesuai")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	user.Password = string(hashed)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("gagal mengupdate password")
	}

	return nil
}

// ── Login ──────────────────────────────────────────────────

func (s *UserService) Login(req models.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("email atau password salah")
	}

	return user, nil
}