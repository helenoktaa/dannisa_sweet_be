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

	products := []models.Produk{
		{
			IDProduk:   "DS001",
			NamaProduk: "Bolen Pisang Coklat Isi 10",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS002",
			NamaProduk: "Bolen Pisang Keju Isi 10",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS003",
			NamaProduk: "Bolen Pisang Mix Isi 10",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS004",
			NamaProduk: "Bolen Pisang Coklat Isi 6",
			HargaJual:  12000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS005",
			NamaProduk: "Bolen Pisang Keju Isi 6",
			HargaJual:  12000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS006",
			NamaProduk: "Bolen Pisang Mix Isi 6",
			HargaJual:  12000,
			Stok:       15,
			IDKategori: "KDS001",
		},
		{
			IDProduk:   "DS007",
			NamaProduk: "Nastar Klasik 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS008",
			NamaProduk: "Nastar Keju 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS009",
			NamaProduk: "Nastar Daun 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS010",
			NamaProduk: "Nastar Roll 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS011",
			NamaProduk: "Choco Roll Tart 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS012",
			NamaProduk: "Kastengels 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS013",
			NamaProduk: "Nutella Coklat 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS014",
			NamaProduk: "Putri Salju 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS015",
			NamaProduk: "Sagu Keju 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS016",
			NamaProduk: "Semprit Jadul 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS017",
			NamaProduk: "Tumbrint Coklat 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS018",
			NamaProduk: "Tumbrint Strawberry 250g",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS002",
		},
		{
			IDProduk:   "DS019",
			NamaProduk: "Bolu Tape",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS003",
		},
		{
			IDProduk:   "DS020",
			NamaProduk: "Bolu Caramel",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS003",
		},
		{
			IDProduk:   "DS021",
			NamaProduk: "Bolu Ketan",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS003",
		},
		{
			IDProduk:   "DS022",
			NamaProduk: "Brownies Lumer Kacang Meses",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS004",
		},
		{
			IDProduk:   "DS023",
			NamaProduk: "Brownies Lumer Coklat Keju",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS004",
		},
		{
			IDProduk:   "DS024",
			NamaProduk: "Brownies Lumer Coklat Almond",
			HargaJual:  20000,
			Stok:       15,
			IDKategori: "KDS004",
		},
		{
			IDProduk:   "DS025",
			NamaProduk: "Brownies Sekat Mix Topping",
			HargaJual:  25000,
			Stok:       15,
			IDKategori: "KDS004",
		},
		{
			IDProduk:   "DS026",
			NamaProduk: "Roti Pisang Coklat",
			HargaJual:  6000,
			Stok:       15,
			IDKategori: "KDS005",
		},
		{
			IDProduk:   "DS027",
			NamaProduk: "Roti Pizza",
			HargaJual:  6000,
			Stok:       15,
			IDKategori: "KDS005",
		},
		{
			IDProduk:   "DS028",
			NamaProduk: "Roti Abon",
			HargaJual:  6000,
			Stok:       15,
			IDKategori: "KDS005",
		},
		{
			IDProduk:   "DS029",
			NamaProduk: "Roti Coklat Keju",
			HargaJual:  6000,
			Stok:       15,
			IDKategori: "KDS005",
		},
		{
			IDProduk:   "DS030",
			NamaProduk: "Roti Meses",
			HargaJual:  5000,
			Stok:       15,
			IDKategori: "KDS005",
		},
	}

	for _, p := range products {
		config.DB.Create(&p)
	}

	log.Printf("Seed berhasil: %d produk ditambahkan", len(products))
}