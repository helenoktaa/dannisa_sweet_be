package repositories

import (
	"fmt"

	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"gorm.io/gorm"
)

type UserMenuRepository struct{}

func NewUserMenuRepository() *UserMenuRepository {
	return &UserMenuRepository{}
}

// ── Generate ID ────────────────────────────────────────────
// Format: UMN01, UMN02, dst.

func (r *UserMenuRepository) GetLastNumber() (int, error) {
	var lastID string
	result := config.DB.Model(&models.UserMenu{}).
		Select("id_user_menu").
		Order("id_user_menu DESC").
		Limit(1).
		Pluck("id_user_menu", &lastID)

	if result.Error != nil || lastID == "" {
		return 0, nil
	}

	var number int
	fmt.Sscanf(lastID, "UMN%d", &number)
	return number, nil
}

// ── Find By User ID ────────────────────────────────────────

func (r *UserMenuRepository) FindByUserID(idUser string) ([]models.UserMenu, error) {
	var menus []models.UserMenu
	result := config.DB.
		Where("id_user = ?", idUser).
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}

// ── Replace (Delete lama + Insert baru) ────────────────────
// Dipanggil saat create user baru atau update menu permissions

func (r *UserMenuRepository) Replace(idUser string, menuKeys []string) error {
	validKeys := models.FilterValidMenuKeys(menuKeys)

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Hapus semua menu lama milik user ini
		if err := tx.
			Where("id_user = ?", idUser).
			Delete(&models.UserMenu{}).Error; err != nil {
			return err
		}

		// Jika tidak ada menu yang dipilih, selesai
		if len(validKeys) == 0 {
			return nil
		}

		// Ambil last number di dalam transaksi
		var lastID string
		tx.Model(&models.UserMenu{}).
			Select("id_user_menu").
			Order("id_user_menu DESC").
			Limit(1).
			Pluck("id_user_menu", &lastID)

		var lastNumber int
		fmt.Sscanf(lastID, "UMN%d", &lastNumber)

		// Build slice UserMenu
		menus := make([]models.UserMenu, 0, len(validKeys))
		for i, key := range validKeys {
			menus = append(menus, models.UserMenu{
				IDUserMenu: fmt.Sprintf("UMN%02d", lastNumber+i+1),
				IDUser:     idUser,
				MenuKey:    key,
			})
		}

		return tx.Create(&menus).Error
	})
}

// ── Delete By User ID ──────────────────────────────────────
// Dipanggil sebelum delete user agar tidak foreign key error

func (r *UserMenuRepository) DeleteByUserID(idUser string) error {
	return config.DB.
		Where("id_user = ?", idUser).
		Delete(&models.UserMenu{}).Error
}

// ── Has Access ─────────────────────────────────────────────
// Cek apakah user punya akses ke menu tertentu

func (r *UserMenuRepository) HasAccess(idUser string, menuKey string) (bool, error) {
	var count int64
	result := config.DB.Model(&models.UserMenu{}).
		Where("id_user = ? AND menu_key = ?", idUser, menuKey).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}