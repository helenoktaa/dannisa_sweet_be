package services

import (
	"errors"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
)

type UserMenuService struct {
	userMenuRepo *repositories.UserMenuRepository
}

func NewUserMenuService(userMenuRepo *repositories.UserMenuRepository) *UserMenuService {
	return &UserMenuService{userMenuRepo: userMenuRepo}
}

// ── Get Menu By User ID ────────────────────────────────────

func (s *UserMenuService) GetByUserID(idUser string) ([]string, error) {
	menus, err := s.userMenuRepo.FindByUserID(idUser)
	if err != nil {
		return nil, errors.New("gagal mengambil data menu")
	}
	return models.ExtractMenuKeys(menus), nil
}

// ── Replace Menu ───────────────────────────────────────────
// Hapus semua menu lama lalu insert yang baru

func (s *UserMenuService) Replace(idUser string, menuKeys []string) error {
	if err := s.userMenuRepo.Replace(idUser, menuKeys); err != nil {
		return errors.New("gagal mengupdate akses menu")
	}
	return nil
}

// ── Delete By User ID ──────────────────────────────────────
// Dipanggil saat user dihapus

func (s *UserMenuService) DeleteByUserID(idUser string) error {
	if err := s.userMenuRepo.DeleteByUserID(idUser); err != nil {
		return errors.New("gagal menghapus akses menu")
	}
	return nil
}

// ── Has Access ─────────────────────────────────────────────
// Dipanggil dari middleware untuk cek akses per menu

func (s *UserMenuService) HasAccess(idUser string, menuKey string) (bool, error) {
	if !models.IsValidMenuKey(menuKey) {
		return false, errors.New("menu key tidak valid")
	}

	ok, err := s.userMenuRepo.HasAccess(idUser, menuKey)
	if err != nil {
		return false, errors.New("gagal mengecek akses menu")
	}
	return ok, nil
}

// ── Get All Available Menus ────────────────────────────────
// Kembalikan semua menu yang tersedia di aplikasi

func (s *UserMenuService) GetAllAvailableMenus() []string {
	return models.AllMenuKeys
}