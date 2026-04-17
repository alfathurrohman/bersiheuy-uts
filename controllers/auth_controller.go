package controllers

import (
        "net/http"

        "github.com/gin-contrib/sessions"
        "github.com/gin-gonic/gin"
        "golang.org/x/crypto/bcrypt"

        "github.com/alfathurrohman/bersiheuy/config"
        "github.com/alfathurrohman/bersiheuy/models"
)

// Halaman Register
func RegisterPage(c *gin.Context) {
        c.HTML(http.StatusOK, "register.html", gin.H{"title": "Daftar - Bersih Euy"})
}

// Halaman Login
func LoginPage(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", gin.H{"title": "Masuk - Bersih Euy"})
}

// Proses Pendaftaran Akun Warga
func Register(c *gin.Context) {
        name := c.PostForm("name")
        email := c.PostForm("email")
        password := c.PostForm("password")

        // Mengacak (Hash) password agar aman di database
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

        user := models.User{
                Name:     name,
                Email:    email,
                Password: string(hashedPassword),
                Role:     "warga", // Semua pendaftar baru otomatis menjadi warga
        }

        config.DB.Create(&user)
        c.Redirect(http.StatusFound, "/login") // Arahkan ke halaman login setelah sukses daftar
}

// Proses Cek Login
func Login(c *gin.Context) {
        email := c.PostForm("email")
        password := c.PostForm("password")

        var user models.User
        // Cari user berdasarkan email
        if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak terdaftar!"})
                return
        }

        // Cek kecocokan password
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
                return
        }

        // Simpan sesi login ke dalam memori aplikasi
        session := sessions.Default(c)
        session.Set("user_id", user.ID)
        session.Set("role", user.Role)
        session.Save()

        // Pisahkan arah halaman berdasarkan Role
        if user.Role == "admin" {
                c.Redirect(http.StatusFound, "/admin")
        } else {
                c.Redirect(http.StatusFound, "/dashboard") // Nanti kita buat halaman dashboard khusus warga
        }
}

// Proses Keluar
func Logout(c *gin.Context) {
        session := sessions.Default(c)
        session.Clear()
        session.Save()
        c.Redirect(http.StatusFound, "/login")
}