package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/alfathurrohman/bersiheuy/config"
	"github.com/alfathurrohman/bersiheuy/controllers"
	"github.com/alfathurrohman/bersiheuy/middlewares"
	"github.com/alfathurrohman/bersiheuy/models"
)

func main() {
	// 1. Inisialisasi Database
	config.ConnectDB()
	err := config.DB.AutoMigrate(&models.User{}, &models.Report{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi:", err)
	}

	// 2. Inisialisasi Aplikasi Gin
	r := gin.Default()

	// 3. Setup Session
	store := cookie.NewStore([]byte("rahasia-negara-bersiheuy"))
	r.Use(sessions.Sessions("bersiheuy_session", store))

	// 4. Load Template HTML
	r.LoadHTMLGlob("templates/*")

	// ==========================================
	// DAFTAR JALUR APLIKASI (ROUTER)
	// ==========================================

	// RUTE PUBLIK (Bebas Akses)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "Selamat Datang - Bersih Euy"})
	})
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.RegisterPage)
	r.POST("/register", controllers.Register)
	r.GET("/logout", controllers.Logout)

	// RUTE WARGA (Wajib Login sebagai apapun)
	warga := r.Group("/")
	warga.Use(middlewares.RequireAuth())
	{
		warga.GET("/dashboard", controllers.UserDashboard)
		warga.POST("/api/reports", controllers.SubmitReport)
		warga.POST("/api/reports/delete/:id", controllers.DeleteReport)
	}

	// RUTE ADMIN (Wajib Login & Wajib Role Admin)
	admin := r.Group("/admin")
	admin.Use(middlewares.RequireAuth(), middlewares.RequireAdmin())
	{
		admin.GET("", controllers.AdminDashboard)
		admin.POST("/update/:id", controllers.UpdateReportStatus)
		admin.POST("/delete/:id", controllers.AdminDeleteReport) // <--- TAMBAHKAN BARIS INI
	}

	// 5. Jalankan Server
	log.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}