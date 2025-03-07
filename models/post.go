package models

import "time"

type Post struct {
	PostID      uint      `gorm:"primaryKey" json:"post_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Content     string    `json:"content"`
	// ImageURL    string    `json:"image_url"`
	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	User     User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes"`
	BookMarks    []BookMark    `gorm:"foreignKey:PostID" json:"bookmarks"`

}
