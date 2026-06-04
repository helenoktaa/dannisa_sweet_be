package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helenoktaa/dannisa_sweet_be/models"
	"github.com/helenoktaa/dannisa_sweet_be/repositories"
	"github.com/helenoktaa/dannisa_sweet_be/services"
)

// inisialisasi sekali — tidak perlu dibuat ulang tiap request
var userMenuService = services.NewUserMenuService(
	repositories.NewUserMenuRepository(),
)


func MenuAccess(menuKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idUserRaw, exists := c.Get("id_user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "User tidak terautentikasi",
				"error_code": "UNAUTHORIZED",
			})
			return
		}

		idUser, ok := idUserRaw.(string)
		if !ok || idUser == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "ID user tidak valid",
				"error_code": "INVALID_USER_ID",
			})
			return
		}

		// ✅ Admin bypass — skip pengecekan menu_keys
		role, _ := c.Get("role")
		if role == "Admin" {
			c.Next()
			return
		}

		// Validasi menu key
		if !models.IsValidMenuKey(menuKey) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"message":    "Menu key tidak dikenali",
				"error_code": "INVALID_MENU_KEY",
			})
			return
		}

		// Cek akses ke DB
		hasAccess, err := userMenuService.HasAccess(idUser, menuKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"message":    "Gagal mengecek akses menu",
				"error_code": "MENU_CHECK_ERROR",
			})
			return
		}

		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "Anda tidak memiliki akses ke menu ini",
				"error_code": "MENU_FORBIDDEN",
			})
			return
		}

		c.Next()
	}
}
