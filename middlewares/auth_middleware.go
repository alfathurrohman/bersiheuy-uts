package middlewares

import (
        "net/http"

        "github.com/gin-contrib/sessions"
        "github.com/gin-gonic/gin"
)

// Wajib Login untuk Warga
func RequireAuth() gin.HandlerFunc {
        return func(c *gin.Context) {
                session := sessions.Default(c)
                userID := session.Get("user_id")

                if userID == nil {
                        // Jika belum login, tendang ke halaman login
                        c.Redirect(http.StatusFound, "/login")
                        c.Abort()
                        return
                }

                // Loloskan dan simpan data ke konteks request
                c.Set("user_id", userID)
                c.Set("role", session.Get("role"))
                c.Next()
        }
}

// Wajib Login untuk Admin
func RequireAdmin() gin.HandlerFunc {
        return func(c *gin.Context) {
                role := c.GetString("role")

                if role != "admin" {
                        // Jika warga mencoba masuk ke halaman admin, tolak!
                        c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! Halaman ini khusus Admin."})
                        c.Abort()
                        return
                }
                c.Next()
        }
}