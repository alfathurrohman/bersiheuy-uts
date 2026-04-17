package config

import (
        "fmt"
        "log"
        "os"

        "github.com/joho/godotenv"
        "gorm.io/driver/mysql"
        "gorm.io/gorm"
)

// DB adalah variabel global agar bisa dipakai di file lain (seperti controller)
var DB *gorm.DB

func ConnectDB() {
        // Memuat variabel dari file .env
        err := godotenv.Load()
        if err != nil {
                log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan sistem")
        }

        // Mengambil data kredensial dari .env
        dbUser := os.Getenv("DB_USER")
        dbPass := os.Getenv("DB_PASSWORD")
        dbHost := os.Getenv("DB_HOST")
        dbPort := os.Getenv("DB_PORT")
        dbName := os.Getenv("DB_NAME")

        // Menyusun string koneksi (DSN - Data Source Name) untuk MySQL
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
                dbUser, dbPass, dbHost, dbPort, dbName)

        // Membuka koneksi menggunakan GORM
        database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
                log.Fatal("Gagal koneksi ke database MySQL:", err)
        }

        DB = database
        fmt.Println("Koneksi ke database MySQL berhasil!")
}