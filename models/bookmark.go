package models

import "time"

type BookMark struct {
	BookMarkID  uint      `gorm:"primaryKey" json:"bookmark_id"`
	PostID      uint      `gorm:"index" json:"post_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	Post Post `gorm:"foreignKey:PostID;references:PostID;constraint:OnDelete:CASCADE;" json:"post"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}

// Explicitly set the table name to match the database table
func (BookMark) TableName() string {
	return "bookmarks"
}

