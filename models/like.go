package models

import "time"

type Like struct {
	LikeID      uint      `gorm:"primaryKey" json:"like_id"`
	PostID      uint      `gorm:"index" json:"post_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	// CommentID   uint      `gorm:"index" json:"comment_id"`
	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	Post Post `gorm:"foreignKey:PostID;references:PostID;constraint:OnDelete:CASCADE;" json:"post"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}
