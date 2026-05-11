package main

import (
	"log"
	"os"

	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/pkg/logger"
	"github.com/helenoktaa/dannisa_sweet_be/routes"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables dari .env file
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 2. Inisialisasi logger (harus pertama agar semua komponen bisa log)
	logger.Init()

	// 3. Inisialisasi Firebase Admin SDK Dannisa Sweet
	// Pastikan file firebase-service-account.json sudah diganti
	// dengan yang baru dari Firebase project Dannisa Sweet
	config.InitFirebase()

	// 4. Inisialisasi koneksi database MySQL (XAMPP)
	config.InitDatabase()

	// 5. AutoMigrate — buat/update tabel otomatis sesuai struct model
	// URUTAN PENTING: tabel yang direferensi FK harus di-migrate duluan
	err := config.DB.AutoMigrate(
		&models.Jabatan{},         // 1. Jabatan dulu (FK dari Users)
		&models.User{},            // 2. Users (FK dari Transaksi)
		&models.Kategori{},        // 3. Kategori dulu (FK dari Produk)
		&models.Produk{},          // 4. Produk (FK dari DetailTransaksi)
		&models.Transaksi{},       // 5. Transaksi (FK dari DetailTransaksi)
		&models.DetailTransaksi{}, // 6. DetailTransaksi paling terakhir
	)
	if err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}
	logger.L.Info("AutoMigrate berhasil — semua tabel sudah siap")

	// 6. Seed data awal jabatan jika tabel masih kosong
	seedJabatan()

	// 7. Setup Gin router dengan semua routes
	router := routes.SetupRouter()

	// 8. Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logger.L.Info("server starting",
		"app", "Dannisa Sweet POS",
		"url", "http://localhost:"+port,
		"health", "http://localhost:"+port+"/v1/health",
	)

	if err := router.Run(":" + port); err != nil {
		logger.L.Error("server gagal berjalan", "error", err)
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

// seedJabatan mengisi data awal jabatan jika tabel jabatan masih kosong
// Dijalankan otomatis saat server pertama kali start
func seedJabatan() {
	var count int64
	config.DB.Model(&models.Jabatan{}).Count(&count)
	if count > 0 {
		// Data jabatan sudah ada, skip seed
		return
	}

	jabatans := []models.Jabatan{
		{NamaJabatan: "Admin", Gaji: 5000000},
		{NamaJabatan: "Kasir", Gaji: 3000000},
	}

	if err := config.DB.Create(&jabatans).Error; err != nil {
		logger.L.Error("Seed jabatan gagal", "error", err)
		return
	}
	logger.L.Info("Seed jabatan berhasil — data Admin dan Kasir sudah ditambahkan")
}