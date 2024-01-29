package models

import (
	"time"
)

type Picture struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int
	User      User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
