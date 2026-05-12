package config

import (
	"fmt"
	"log"
	"os"

	"github.com/helenoktaa/dannisa_sweet_be/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB adalah instance GORM global yang dipakai di seluruh aplikasi
var DB *gorm.DB

func InitDatabase() {
	// Ambil konfigurasi dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format DSN untuk MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	// Konfigurasi GORM
	gormConfig := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
        SingularTable: true,
    },
	}

	// Buka koneksi
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// Setup connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	// AutoMigrate: urutan penting!
	// Tabel yang tidak punya FK harus dibuat duluan
	err = DB.AutoMigrate(
		&models.Jabatan{},         // 1. Jabatan dulu (tidak ada FK)
		&models.User{},            // 2. User → FK ke Jabatan
		&models.Kategori{},        // 3. Kategori (tidak ada FK)
		&models.Produk{},          // 4. Produk → FK ke Kategori
		&models.Transaksi{},       // 5. Transaksi → FK ke User
		&models.DetailTransaksi{}, // 6. DetailTransaksi → FK ke Transaksi & Produk
	)
	if err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}

	log.Println("✅ Database dannisa_sweet terhubung dan tabel sudah di-migrate")
}
