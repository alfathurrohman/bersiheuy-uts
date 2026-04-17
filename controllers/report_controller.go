package controllers

import (
        "context"
        "fmt"
        "net/http"
        "os"
        "time"

        "github.com/aws/aws-sdk-go-v2/aws"
        awsConfig "github.com/aws/aws-sdk-go-v2/config"
        "github.com/aws/aws-sdk-go-v2/service/s3"
        "github.com/gin-contrib/sessions"
        "github.com/gin-gonic/gin"

        "github.com/alfathurrohman/bersiheuy/config"
        "github.com/alfathurrohman/bersiheuy/models"
)

// ==========================================
// FITUR UNTUK WARGA PELAPOR
// ==========================================

// Halaman Dashboard Warga (Melihat riwayat laporannya sendiri)
func UserDashboard(c *gin.Context) {
        session := sessions.Default(c)
        userID := session.Get("user_id")

        var reports []models.Report
        // Cari laporan di database yang HANYA milik user yang sedang login
        config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&reports)

        c.HTML(http.StatusOK, "dashboard.html", gin.H{
                "title":   "Dashboard - Bersih Euy",
                "reports": reports,
        })
}

// Proses Mengirim Laporan Baru
func SubmitReport(c *gin.Context) {
        session := sessions.Default(c)
        userID := session.Get("user_id").(uint) // Identifikasi siapa yang melapor

        location := c.PostForm("location")
        description := c.PostForm("description")

        file, header, err := c.Request.FormFile("photo")
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Foto wajib diunggah!"})
                return
        }
        defer file.Close()

        // 1. Upload ke S3 AWS
        cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(os.Getenv("AWS_REGION")))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat konfigurasi AWS"})
                return
        }
        s3Client := s3.NewFromConfig(cfg)
        bucketName := os.Getenv("AWS_S3_BUCKET")
        fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

        _, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
                Bucket:      aws.String(bucketName),
                Key:         aws.String(fileName),
                Body:        file,
                ContentType: aws.String(header.Header.Get("Content-Type")), // Agar foto tidak otomatis ter-download
        })
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload foto: " + err.Error()})
                return
        }

        photoURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("AWS_REGION"), fileName)

        // 2. Simpan ke Database
        report := models.Report{
                UserID:      userID, // Menempelkan ID warga pada laporan
                Location:    location,
                Description: description,
                PhotoURL:    photoURL,
                Status:      "Menunggu",
        }
        config.DB.Create(&report)

        // Arahkan kembali ke dashboard warga
        c.Redirect(http.StatusFound, "/dashboard")
}

// Fitur Hapus Laporan (Hanya bisa jika statusnya masih "Menunggu")
func DeleteReport(c *gin.Context) {
        id := c.Param("id") // Ambil ID laporan dari URL
        session := sessions.Default(c)
        userID := session.Get("user_id")

        var report models.Report
        // Pastikan laporan itu ada dan BENAR milik user tersebut
        if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&report).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Laporan tidak ditemukan atau Anda tidak berhak menghapusnya"})
                return
        }

        // Kunci keamanan: Kalau sudah diproses admin, tidak boleh dihapus!
        if report.Status != "Menunggu" {
                c.JSON(http.StatusForbidden, gin.H{"error": "Laporan yang sudah diproses admin tidak dapat dihapus!"})
                return
        }

        config.DB.Delete(&report)
        c.Redirect(http.StatusFound, "/dashboard")
}

// ==========================================
// FITUR KHUSUS ADMIN
// ==========================================

// Halaman Dashboard Admin (Melihat SEMUA laporan)
func AdminDashboard(c *gin.Context) {
        var reports []models.Report
        // Ambil SEMUA laporan, dan gunakan "Preload" untuk mengambil nama warga pelapornya
        config.DB.Preload("User").Order("created_at desc").Find(&reports)

        c.HTML(http.StatusOK, "admin.html", gin.H{
                "title":   "Admin Area - Bersih Euy",
                "reports": reports,
        })
}

// Fitur Admin untuk mengubah status dan membalas pesan
func UpdateReportStatus(c *gin.Context) {
        id := c.Param("id")
        status := c.PostForm("status")
        adminReply := c.PostForm("admin_reply")

        var report models.Report
        if err := config.DB.First(&report, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Laporan tidak ditemukan"})
                return
        }

        // Update status dan pesan balasannya
        report.Status = status
        report.AdminReply = adminReply
        config.DB.Save(&report)

        c.Redirect(http.StatusFound, "/admin")
}

// Fitur Admin untuk menghapus laporan secara permanen
func AdminDeleteReport(c *gin.Context) {
        id := c.Param("id")

        var report models.Report
        // Cari laporan berdasarkan ID
        if err := config.DB.First(&report, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Laporan tidak ditemukan"})
                return
        }

        // Hapus laporan dari database
        config.DB.Delete(&report)

        // Arahkan kembali ke halaman admin
        c.Redirect(http.StatusFound, "/admin")
}