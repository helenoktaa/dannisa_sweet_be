package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/handlers"
	"github.com/helenoktaa/dannisa_sweet_be/middleware"
	"github.com/helenoktaa/dannisa_sweet_be/models"
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
	authHandler        := handlers.NewAuthHandler()
	produkHandler      := handlers.NewProductHandler()
	kategoriHandler    := handlers.NewKategoriHandler()
	jabatanHandler     := handlers.NewJabatanHandler()
	transaksiHandler   := handlers.NewTransaksiHandler()
	userHandler        := handlers.NewUserHandler()
	stokHistoryHandler := handlers.NewStokHistoryHandler()
	dashboardHandler   := handlers.NewDashboardHandler()
	markdownHandler    := handlers.NewMarkdownPricingHandler()

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

		// ── Auth routes (public) ──────────────────────────────
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// ── Protected routes ──────────────────────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// ── Profile ───────────────────────────────────────
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.PUT("/auth/password", authHandler.UpdatePassword)

			// ── Produk ────────────────────────────────────────
			produk := protected.Group("/produk")
			produk.Use(middleware.MenuAccess(models.MenuProduk))
			{
				produk.GET("", produkHandler.GetAll)
				produk.GET("/:id", produkHandler.GetByID)

				adminProduk := produk.Group("")
				adminProduk.Use(middleware.AdminOnly())
				{
					adminProduk.POST("", produkHandler.Create)
					adminProduk.PUT("/:id", produkHandler.Update)
					adminProduk.DELETE("/:id", produkHandler.Delete)
				}
			}

			// ── Markdown Pricing (Admin only) ─────────────────
			markdown := protected.Group("/markdown")
			markdown.Use(middleware.AdminOnly())
			{
				markdown.POST("", markdownHandler.SetMarkdown)
				markdown.GET("/:id_produk", markdownHandler.GetMarkdown)
				markdown.PATCH("/:id_produk/override", markdownHandler.OverrideManual)
				markdown.DELETE("/:id_produk/override", markdownHandler.HapusOverrideManual)
			}

			// ── Kategori ──────────────────────────────────────
			kategori := protected.Group("/kategori")
			{
				kategori.GET("", kategoriHandler.GetAll)
				kategori.GET("/:id", kategoriHandler.GetByID)

				adminKategori := kategori.Group("")
				adminKategori.Use(middleware.AdminOnly())
				{
					adminKategori.POST("", kategoriHandler.Create)
					adminKategori.PUT("/:id", kategoriHandler.Update)
					adminKategori.DELETE("/:id", kategoriHandler.Delete)
				}
			}

			// ── Jabatan ───────────────────────────────────────
			jabatan := protected.Group("/jabatan")
			{
				jabatan.GET("", jabatanHandler.GetAll)
				jabatan.GET("/:id", jabatanHandler.GetByID)

				adminJabatan := jabatan.Group("")
				adminJabatan.Use(middleware.AdminOnly())
				{
					adminJabatan.POST("", jabatanHandler.Create)
					adminJabatan.PUT("/:id", jabatanHandler.Update)
					adminJabatan.DELETE("/:id", jabatanHandler.Delete)
				}
			}

			// ── Transaksi ─────────────────────────────────────
			transaksi := protected.Group("/transaksi")
			transaksi.Use(middleware.MenuAccess(models.MenuTransaksi))
			{
				transaksi.POST("", transaksiHandler.Create)
				transaksi.GET("", transaksiHandler.GetAll)
				transaksi.GET("/pre-order/aktif", transaksiHandler.GetPreOrderAktif)
				transaksi.PATCH("/:id/status-order", transaksiHandler.UpdateStatusOrder)
				transaksi.GET("/:id", transaksiHandler.GetByID)
				transaksi.GET("/:id/invoice", transaksiHandler.GetInvoice)

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
				users.GET("", userHandler.GetAll)
				users.POST("", userHandler.Create)
				users.GET("/:id", userHandler.GetByID)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.PUT("/:id/password", userHandler.UpdatePassword)

				// Menu permissions
				users.GET("/:id/menus", userHandler.GetUserMenus)
				users.PUT("/:id/menus", userHandler.UpdateUserMenus)
			}

			// ── Available Menus (Admin only) ──────────────────
			protected.GET("/menus", middleware.AdminOnly(), userHandler.GetAvailableMenus)

			// ── Stok History (Admin only) ──────────────────────
			stokHistory := protected.Group("/stok-history")
			stokHistory.Use(middleware.AdminOnly())
			{
				stokHistory.POST("", stokHistoryHandler.Create)
				stokHistory.GET("", stokHistoryHandler.GetAll)
			}

			// ── Dashboard ─────────────────────────────────────
			dashboard := protected.Group("/dashboard")
			dashboard.Use(middleware.MenuAccess(models.MenuDashboard))
			{
				dashboard.GET("", dashboardHandler.GetDashboard)
				dashboard.GET("/harian", dashboardHandler.GetDashboardHarian)
			}

			// ── Laporan (via transaksi, Admin only) ───────────
			protected.GET("/laporan",
				middleware.AdminOnly(),
				middleware.MenuAccess(models.MenuLaporan),
				transaksiHandler.GetLaporan,
			)
		}
	}

	return r
}