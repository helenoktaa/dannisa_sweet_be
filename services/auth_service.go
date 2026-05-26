package services

import (
	"fmt"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
	}
}


// REGISTER

func (s *AuthService) Register(req models.RegisterRequest) (*models.UserResponse, error) {

	// cek email sudah digunakan atau belum
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	 // ── Generate id_user otomatis ──────────────────────────
    lastNumber, err := s.userRepo.GetLastNumber()
    if err != nil {
        return nil, errors.New("gagal generate ID user")
    }
    generatedID := fmt.Sprintf("UDS%02d", lastNumber+1)

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, errors.New("gagal hash password")
	}

	// buat object user
	user := models.User{
		 IDUser:       generatedID,
		NamaUser:      req.NamaUser,
		Email:         req.Email,
		Password:      string(hashedPassword),
		RekPembayaran: req.RekPembayaran,
		Whatsapp:      req.Whatsapp,
		IDJabatan:     req.IDJabatan,
	}

	// simpan ke database
	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("gagal membuat akun")
	}

	// ambil data lengkap beserta relasi jabatan
	createdUser, err := s.userRepo.FindByID(user.IDUser)
	if err != nil {
		return nil, err
	}

	// response
	response := &models.UserResponse{
		IDUser:        createdUser.IDUser,
		NamaUser:      createdUser.NamaUser,
		Email:         createdUser.Email,
		RekPembayaran: createdUser.RekPembayaran,
		Whatsapp:      createdUser.Whatsapp,
		Jabatan: models.JabatanResponse{
			IDJabatan:   createdUser.Jabatan.IDJabatan,
			NamaJabatan: createdUser.Jabatan.NamaJabatan,
			Gaji:        createdUser.Jabatan.Gaji,
		},
	}

	return response, nil
}


// LOGIN

func (s *AuthService) Login(email string, password string) (*models.LoginResponse, error) {

	// cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// generate JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	// response
	response := &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
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
		},
	}

	return response, nil
}


// GET PROFILE

func (s *AuthService) GetProfile(idUser string) (*models.UserResponse, error) {

	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	response := &models.UserResponse{
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
	}

	return response, nil
}


// UPDATE PROFILE

func (s *AuthService) UpdateProfile(
	idUser string,
	req models.UpdateUserRequest,
) (*models.UserResponse, error) {

	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	user.NamaUser = req.NamaUser
	user.RekPembayaran = req.RekPembayaran
	user.Whatsapp = req.Whatsapp

	if req.IDJabatan != "" {
		user.IDJabatan = req.IDJabatan
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("gagal update profil")
	}

	updatedUser, _ := s.userRepo.FindByID(idUser)

	response := &models.UserResponse{
		IDUser:        updatedUser.IDUser,
		NamaUser:      updatedUser.NamaUser,
		Email:         updatedUser.Email,
		RekPembayaran: updatedUser.RekPembayaran,
		Whatsapp:      updatedUser.Whatsapp,
		Jabatan: models.JabatanResponse{
			IDJabatan:   updatedUser.Jabatan.IDJabatan,
			NamaJabatan: updatedUser.Jabatan.NamaJabatan,
			Gaji:        updatedUser.Jabatan.Gaji,
		},
	}

	return response, nil
}


// UPDATE PASSWORD

func (s *AuthService) UpdatePassword(
	idUser string,
	req models.UpdatePasswordRequest,
) error {

	user, err := s.userRepo.FindByID(idUser)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// cek password lama
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.PasswordLama),
	)

	if err != nil {
		return errors.New("password lama salah")
	}

	// hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.PasswordBaru),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return errors.New("gagal hash password")
	}

	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("gagal update password")
	}

	return nil
}


// GENERATE JWT

func (s *AuthService) generateJWT(user *models.User) (string, error) {

	claims := jwt.MapClaims{
		"id_user": user.IDUser,
		"email":   user.Email,
		"role":    user.Jabatan.NamaJabatan,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}