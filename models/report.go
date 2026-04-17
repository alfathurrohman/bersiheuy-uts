package models

import "time"

type Report struct {
        ID          uint      `gorm:"primaryKey" json:"id"`
        UserID      uint      `json:"user_id"` // Menyambungkan laporan ke user tertentu
        Location    string    `gorm:"type:text;not null" json:"location"`
        Description string    `gorm:"type:text;not null" json:"description"`
        PhotoURL    string    `gorm:"type:text;not null" json:"photo_url"`
        Status      string    `gorm:"type:varchar(50);default:'Menunggu'" json:"status"` // Menunggu, Disetujui, Ditolak
        AdminReply  string    `gorm:"type:text" json:"admin_reply"` // Tempat admin membalas pesan
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`

        // Relasi: Satu Laporan dimiliki oleh satu User
        User        User      `gorm:"foreignKey:UserID" json:"user"`
}