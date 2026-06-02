package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/handlers"
	"github.com/helenoktaa/dannisa_sweet_be/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.HTTPLogger())

	// ─── CORS Middleware ───────────────────────────────────────
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ─── Init handlers ────────────────────────────────────────
	authHandler := handlers.NewAuthHandler()
	produkHandler := handlers.NewProductHandler()
	kategoriHandler := handlers.NewKategoriHandler()
	jabatanHandler := handlers.NewJabatanHandler()
	transaksiHandler := handlers.NewTransaksiHandler()
	userHandler := handlers.NewUserHandler()
	stokHistoryHandler := handlers.NewStokHistoryHandler()
	dashboardHandler := handlers.NewDashboardHandler()

	// ─── API v1 group ─────────────────────────────────────────
	v1 := r.Group("/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "dannisa-sweet-backend",
			})
		})

		// ── Auth routes (public, tidak butuh JWT) ─────────────
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)       // POST /v1/auth/login
			auth.POST("/register", authHandler.Register) // POST /v1/auth/register
		}

		// ── Protected routes (semua butuh JWT) ────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// ── Profile (semua role bisa akses) ───────────────
			protected.GET("/auth/profile", authHandler.GetProfile)      // GET /v1/auth/profile
			protected.PUT("/auth/profile", authHandler.UpdateProfile)   // PUT /v1/auth/profile
			protected.PUT("/auth/password", authHandler.UpdatePassword) // PUT /v1/auth/password

			// ── Produk (GET: Admin & Kasir | CUD: Admin only) ─
			produk := protected.Group("/produk")
			{
				produk.GET("", produkHandler.GetAll)      // GET /v1/produk
				produk.GET("/:id", produkHandler.GetByID) // GET /v1/produk/:id

				adminProduk := produk.Group("")
				adminProduk.Use(middleware.AdminOnly())
				{
					adminProduk.POST("", produkHandler.Create)       // POST   /v1/produk
					adminProduk.PUT("/:id", produkHandler.Update)    // PUT    /v1/produk/:id
					adminProduk.DELETE("/:id", produkHandler.Delete) // DELETE /v1/produk/:id
				}
			}

			// ── Kategori (GET: Admin & Kasir | CUD: Admin only) ─
			kategori := protected.Group("/kategori")
			{
				kategori.GET("", kategoriHandler.GetAll)      // GET /v1/kategori
				kategori.GET("/:id", kategoriHandler.GetByID) // GET /v1/kategori/:id

				adminKategori := kategori.Group("")
				adminKategori.Use(middleware.AdminOnly())
				{
					adminKategori.POST("", kategoriHandler.Create)       // POST   /v1/kategori
					adminKategori.PUT("/:id", kategoriHandler.Update)    // PUT    /v1/kategori/:id
					adminKategori.DELETE("/:id", kategoriHandler.Delete) // DELETE /v1/kategori/:id
				}
			}

			// ── Jabatan (GET: Admin & Kasir | CUD: Admin only) ─
			jabatan := protected.Group("/jabatan")
			{
				jabatan.GET("", jabatanHandler.GetAll)      // GET /v1/jabatan
				jabatan.GET("/:id", jabatanHandler.GetByID) // GET /v1/jabatan/:id

				adminJabatan := jabatan.Group("")
				adminJabatan.Use(middleware.AdminOnly())
				{
					adminJabatan.POST("", jabatanHandler.Create)       // POST   /v1/jabatan
					adminJabatan.PUT("/:id", jabatanHandler.Update)    // PUT    /v1/jabatan/:id
					adminJabatan.DELETE("/:id", jabatanHandler.Delete) // DELETE /v1/jabatan/:id
				}
			}

			// ── Transaksi ─────────────────────────────────────────────
			transaksi := protected.Group("/transaksi")
			{
				// Admin & Kasir boleh akses
				transaksi.POST("", transaksiHandler.Create)
				transaksi.GET("", transaksiHandler.GetAll)
				transaksi.GET("/pre-order/aktif", transaksiHandler.GetPreOrderAktif)
				transaksi.PATCH("/:id/status-order", transaksiHandler.UpdateStatusOrder)
				transaksi.GET("/:id", transaksiHandler.GetByID)
				transaksi.GET("/:id/invoice", transaksiHandler.GetInvoice)

				// Admin only
				adminTrx := transaksi.Group("")
				adminTrx.Use(middleware.AdminOnly())

				{
					adminTrx.GET("/laporan", transaksiHandler.GetLaporan)
					adminTrx.PUT("/:id/status", transaksiHandler.UpdateStatus)

				}

			}

			// ── User Management (Admin only) ──────────────────
			users := protected.Group("/users")
			users.Use(middleware.AdminOnly())
			{
				users.GET("", userHandler.GetAll)        // GET    /v1/users
				users.GET("/:id", userHandler.GetByID)   // GET    /v1/users/:id
				users.PUT("/:id", userHandler.Update)    // PUT    /v1/users/:id
				users.DELETE("/:id", userHandler.Delete) // DELETE /v1/users/:id
			}

			// ── Stok History (Admin only) ──────────────────────────────
			stokHistory := protected.Group("/stok-history")
			stokHistory.Use(middleware.AdminOnly())
			{
				stokHistory.POST("", stokHistoryHandler.Create) // catat perubahan stok
				stokHistory.GET("", stokHistoryHandler.GetAll)  // lihat history
			}

			// ── Dashboard (Admin & Kasir) ─────────────────────
			protected.GET("/dashboard", dashboardHandler.GetDashboard)
		}
	}

	return r
}
