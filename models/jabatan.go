package models

// Jabatan menyimpan data jabatan karyawan (Admin / Kasir)
type Jabatan struct {
	IDJabatan int `gorm:"primaryKey;autoIncrement" json:"id_jabatan"`

	NamaJabatan string  `gorm:"not null;size:50" json:"nama_jabatan"`
	Gaji        float64 `gorm:"not null"         json:"gaji"`

	// Relasi: satu jabatan punya banyak user
	Users []User `gorm:"foreignKey:IDJabatan" json:"users,omitempty"`
}

// DTO Request

type CreateJabatanRequest struct {
	NamaJabatan string  `json:"nama_jabatan" binding:"required"`
	Gaji        float64 `json:"gaji"         binding:"required,min=0"`
}

type UpdateJabatanRequest struct {
	NamaJabatan string  `json:"nama_jabatan"`
	Gaji        float64 `json:"gaji" binding:"omitempty,min=0"`
}

// DTO Response

type JabatanResponse struct {
	IDJabatan int `json:"id_jabatan"`

	NamaJabatan string  `json:"nama_jabatan"`
	Gaji        float64 `json:"gaji"`
}