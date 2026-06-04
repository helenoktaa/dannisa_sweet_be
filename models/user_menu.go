package models

import "time"

// ── Konstanta Menu Key ─────────────────────────────────────
const (
	MenuDashboard  = "dashboard"
	MenuTransaksi  = "transaksi"
	MenuProduk     = "produk"
	MenuLaporan    = "laporan"
	MenuKelolaUser = "kelola_user"
)

// AllMenuKeys berisi semua menu yang tersedia di aplikasi
var AllMenuKeys = []string{
	MenuDashboard,
	MenuTransaksi,
	MenuProduk,
	MenuLaporan,
	MenuKelolaUser,
}

// ── Model ──────────────────────────────────────────────────

// UserMenu menyimpan hak akses menu per user
type UserMenu struct {
	IDUserMenu string    `gorm:"primaryKey;size:20"         json:"id_user_menu"`
	IDUser     string    `gorm:"not null;size:20;index"     json:"id_user"`
	MenuKey    string    `gorm:"not null;size:50"           json:"menu_key"`
	CreatedAt  time.Time `gorm:"autoCreateTime"             json:"created_at"`

	// Relasi
	User User `gorm:"foreignKey:IDUser;references:IDUser" json:"-"`
}

// TableName override nama tabel
func (UserMenu) TableName() string {
	return "user_menus"
}

// ── Helper ─────────────────────────────────────────────────

// ExtractMenuKeys mengubah []UserMenu menjadi []string menu key
func ExtractMenuKeys(menus []UserMenu) []string {
	keys := make([]string, 0, len(menus))
	for _, m := range menus {
		keys = append(keys, m.MenuKey)
	}
	return keys
}

// BuildUserMenus membuat slice UserMenu dari idUser dan list menuKey
func BuildUserMenus(idUser string, menuKeys []string, generateID func() string) []UserMenu {
	menus := make([]UserMenu, 0, len(menuKeys))
	for _, key := range menuKeys {
		if IsValidMenuKey(key) {
			menus = append(menus, UserMenu{
				IDUserMenu: generateID(),
				IDUser:     idUser,
				MenuKey:    key,
			})
		}
	}
	return menus
}

// IsValidMenuKey memvalidasi apakah menu key terdaftar
func IsValidMenuKey(key string) bool {
	for _, valid := range AllMenuKeys {
		if key == valid {
			return true
		}
	}
	return false
}

// HasMenuAccess mengecek apakah user punya akses ke menu tertentu
func HasMenuAccess(menus []UserMenu, menuKey string) bool {
	for _, m := range menus {
		if m.MenuKey == menuKey {
			return true
		}
	}
	return false
}

// FilterValidMenuKeys memfilter dan hanya mengembalikan menu key yang valid
func FilterValidMenuKeys(menuKeys []string) []string {
	valid := make([]string, 0, len(menuKeys))
	for _, key := range menuKeys {
		if IsValidMenuKey(key) {
			valid = append(valid, key)
		}
	}
	return valid
}