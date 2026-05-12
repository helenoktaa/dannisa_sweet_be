package models

// User menyimpan data akun karyawan (Admin atau Kasir)
type User struct {
	IDUser        string `gorm:"primaryKey;size:20"        json:"id_user"`
	NamaUser      string `gorm:"not null;size:100"         json:"nama_user"`
	Email         string `gorm:"not null;unique;size:100"  json:"email"`
	Password      string `gorm:"not null"                  json:"-"` // disembunyikan dari response JSON
	RekPembayaran string `gorm:"size:50"                   json:"rek_pembayaran"`
	Whatsapp      string `gorm:"size:20"                   json:"whatsapp"`
	IDJabatan     string `gorm:"not null;size:20;index"         json:"id_jabatan"`

	// Relasi
	Jabatan Jabatan `gorm:"foreignKey:IDJabatan;references:IDJabatan" json:"jabatan,omitempty"`
}

// DTO
type RegisterRequest struct {
	IDUser        string `json:"id_user"         binding:"required"`
	NamaUser      string `json:"nama_user"       binding:"required"`
	Email         string `json:"email"           binding:"required,email"`
	Password      string `json:"password"        binding:"required,min=6"`
	RekPembayaran string `json:"rek_pembayaran"`
	Whatsapp      string `json:"whatsapp"`
	IDJabatan     string `json:"id_jabatan"      binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	NamaUser      string `json:"nama_user"`
	RekPembayaran string `json:"rek_pembayaran"`
	Whatsapp      string `json:"whatsapp"`
	IDJabatan     string `json:"id_jabatan"`
}

type UpdatePasswordRequest struct {
	PasswordLama string `json:"password_lama" binding:"required"`
	PasswordBaru string `json:"password_baru" binding:"required,min=6"`
}

// Response
type UserResponse struct {
	IDUser        string          `json:"id_user"`
	NamaUser      string          `json:"nama_user"`
	Email         string          `json:"email"`
	RekPembayaran string          `json:"rek_pembayaran"`
	Whatsapp      string          `json:"whatsapp"`
	Jabatan       JabatanResponse `json:"jabatan"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
