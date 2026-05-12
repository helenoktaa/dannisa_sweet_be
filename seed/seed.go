package main

import (
	"log"

	"github.com/helenoktaa/dannisa_sweet_be/config"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitDatabase()

	// ── Seed Kategori ──────────────────────────────────────────
	categories := []models.Kategori{
		{IDKategori: "KDS001", NamaKategori: "Bolen"},
		{IDKategori: "KDS002", NamaKategori: "Kue Kering"},
		{IDKategori: "KDS003", NamaKategori: "Bolu"},
		{IDKategori: "KDS004", NamaKategori: "Brownies"},
		{IDKategori: "KDS005", NamaKategori: "Roti"},
	}

	for _, c := range categories {
		config.DB.FirstOrCreate(&c, models.Kategori{IDKategori: c.IDKategori})
	}
	log.Println("✅ Seed kategori selesai")

	// ── Seed Produk ────────────────────────────────────────────
	products := []models.Produk{
		// Bolen
		{IDProduk: "DS001", NamaProduk: "Bolen Pisang Coklat Isi 10", HargaModal: 13000, HargaJual: 20000, Stok: 15, IDKategori: "KDS001"},
		{IDProduk: "DS002", NamaProduk: "Bolen Pisang Keju Isi 10",   HargaModal: 13000, HargaJual: 20000, Stok: 15, IDKategori: "KDS001"},
		{IDProduk: "DS003", NamaProduk: "Bolen Pisang Mix Isi 10",    HargaModal: 13000, HargaJual: 20000, Stok: 15, IDKategori: "KDS001"},
		{IDProduk: "DS004", NamaProduk: "Bolen Pisang Coklat Isi 6",  HargaModal: 8000,  HargaJual: 12000, Stok: 15, IDKategori: "KDS001"},
		{IDProduk: "DS005", NamaProduk: "Bolen Pisang Keju Isi 6",    HargaModal: 8000,  HargaJual: 12000, Stok: 15, IDKategori: "KDS001"},
		{IDProduk: "DS006", NamaProduk: "Bolen Pisang Mix Isi 6",     HargaModal: 8000,  HargaJual: 12000, Stok: 15, IDKategori: "KDS001"},

		// Kue Kering
		{IDProduk: "DS007", NamaProduk: "Nastar Klasik 250g",      HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS008", NamaProduk: "Nastar Keju 250g",        HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS009", NamaProduk: "Nastar Daun 250g",        HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS010", NamaProduk: "Nastar Roll 250g",        HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS011", NamaProduk: "Choco Roll Tart 250g",    HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS012", NamaProduk: "Kastengels 250g",         HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS013", NamaProduk: "Nutella Coklat 250g",     HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS014", NamaProduk: "Putri Salju 250g",        HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS015", NamaProduk: "Sagu Keju 250g",          HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS016", NamaProduk: "Semprit Jadul 250g",      HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS017", NamaProduk: "Tumbrint Coklat 250g",    HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},
		{IDProduk: "DS018", NamaProduk: "Tumbrint Strawberry 250g",HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS002"},

		// Bolu
		{IDProduk: "DS019", NamaProduk: "Bolu Tape",    HargaModal: 13000, HargaJual: 20000, Stok: 15, IDKategori: "KDS003"},
		{IDProduk: "DS020", NamaProduk: "Bolu Caramel", HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS003"},
		{IDProduk: "DS021", NamaProduk: "Bolu Ketan",   HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS003"},

		// Brownies
		{IDProduk: "DS022", NamaProduk: "Brownies Lumer Kacang Meses",  HargaModal: 12000, HargaJual: 20000, Stok: 15, IDKategori: "KDS004"},
		{IDProduk: "DS023", NamaProduk: "Brownies Lumer Coklat Keju",   HargaModal: 12000, HargaJual: 20000, Stok: 15, IDKategori: "KDS004"},
		{IDProduk: "DS024", NamaProduk: "Brownies Lumer Coklat Almond", HargaModal: 12000, HargaJual: 20000, Stok: 15, IDKategori: "KDS004"},
		{IDProduk: "DS025", NamaProduk: "Brownies Sekat Mix Topping",   HargaModal: 15000, HargaJual: 25000, Stok: 15, IDKategori: "KDS004"},

		// Roti
		{IDProduk: "DS026", NamaProduk: "Roti Pisang Coklat", HargaModal: 3500, HargaJual: 6000, Stok: 15, IDKategori: "KDS005"},
		{IDProduk: "DS027", NamaProduk: "Roti Pizza",         HargaModal: 3500, HargaJual: 6000, Stok: 15, IDKategori: "KDS005"},
		{IDProduk: "DS028", NamaProduk: "Roti Abon",          HargaModal: 3500, HargaJual: 6000, Stok: 15, IDKategori: "KDS005"},
		{IDProduk: "DS029", NamaProduk: "Roti Coklat Keju",   HargaModal: 3500, HargaJual: 6000, Stok: 15, IDKategori: "KDS005"},
		{IDProduk: "DS030", NamaProduk: "Roti Meses",         HargaModal: 3000, HargaJual: 5000, Stok: 15, IDKategori: "KDS005"},
	}

	for _, p := range products {
		config.DB.FirstOrCreate(&p, models.Produk{IDProduk: p.IDProduk})
	}

	log.Printf("✅ Seed selesai: %d kategori, %d produk ditambahkan", len(categories), len(products))
}