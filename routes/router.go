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
	authHandler      := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()
	kategoriHandler  := handlers.NewKategoriHandler()
	transaksiHandler := handlers.NewTransaksiHandler()
	userHandler      := handlers.NewUserHandler()

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
			auth.POST("/login",    authHandler.Login)    // POST /v1/auth/login
			auth.POST("/register", authHandler.Register) // POST /v1/auth/register (nanti di-protect admin)
		}

		// ── Protected routes (semua butuh JWT) ────────────────
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// ── Profile (semua role bisa akses) ───────────────
			protected.GET("/auth/profile",    authHandler.GetProfile)    // GET    /v1/auth/profile
			protected.PUT("/auth/profile",    authHandler.UpdateProfile)  // PUT    /v1/auth/profile
			protected.PUT("/auth/password",   authHandler.UpdatePassword) // PUT    /v1/auth/password

			// ── Produk ────────────────────────────────────────
			produk := protected.Group("/produk")
			{
				produk.GET("",     productHandler.GetAll)   // GET /v1/produk         (Admin & Kasir)
				produk.GET("/:id", productHandler.GetByID)  // GET /v1/produk/:id     (Admin & Kasir)

				// Hanya Admin
				adminProduk := produk.Group("")
				adminProduk.Use(middleware.AdminOnly())
				{
					adminProduk.POST("",      productHandler.Create) // POST   /v1/produk
					adminProduk.PUT("/:id",   productHandler.Update) // PUT    /v1/produk/:id
					adminProduk.DELETE("/:id", productHandler.Delete) // DELETE /v1/produk/:id
				}
			}

			// ── Kategori ──────────────────────────────────────
			kategori := protected.Group("/kategori")
			{
				kategori.GET("",     kategoriHandler.GetAll)   // GET /v1/kategori       (Admin & Kasir)
				kategori.GET("/:id", kategoriHandler.GetByID)  // GET /v1/kategori/:id   (Admin & Kasir)

				// Hanya Admin
				adminKategori := kategori.Group("")
				adminKategori.Use(middleware.AdminOnly())
				{
					adminKategori.POST("",       kategoriHandler.Create) // POST   /v1/kategori
					adminKategori.PUT("/:id",    kategoriHandler.Update) // PUT    /v1/kategori/:id
					adminKategori.DELETE("/:id", kategoriHandler.Delete) // DELETE /v1/kategori/:id
				}
			}

			// ── Transaksi ─────────────────────────────────────
			transaksi := protected.Group("/transaksi")
			{
				// Admin & Kasir bisa buat dan lihat transaksi
				transaksi.POST("",     transaksiHandler.Create)   // POST /v1/transaksi       (input transaksi baru)
				transaksi.GET("/:id",  transaksiHandler.GetByID)  // GET  /v1/transaksi/:id   (detail transaksi)
				transaksi.GET("/:id/struk", transaksiHandler.GetStruk) // GET /v1/transaksi/:id/struk (generate struk)

				// Hanya Admin (laporan semua transaksi)
				adminTrx := transaksi.Group("")
				adminTrx.Use(middleware.AdminOnly())
				{
					adminTrx.GET("", transaksiHandler.GetAll) // GET /v1/transaksi (semua transaksi)
				}
			}

			// ── User Management (Admin only) ──────────────────
			users := protected.Group("/users")
			users.Use(middleware.AdminOnly())
			{
				users.GET("",        userHandler.GetAll)    // GET    /v1/users
				users.GET("/:id",    userHandler.GetByID)   // GET    /v1/users/:id
				users.PUT("/:id",    userHandler.Update)    // PUT    /v1/users/:id
				users.DELETE("/:id", userHandler.Delete)    // DELETE /v1/users/:id
			}
		}
	}

	return r
}