package models

import "time"

type Comment struct {
	CommentID   uint      `gorm:"primaryKey" json:"comment_id"`
	PostID      *uint     `gorm:"index" json:"post_id"` // Nullable because a comment can be a reply to another comment.
	UserID      uint      `gorm:"index" json:"user_id"`
	ParentID    *uint     `gorm:"index" json:"parent_id"` // For nested comments (self-referencing)
	Content     string    `json:"content"`
	// LikeCount     int64     `gorm:"default:0" json:"like_count"`
	// BookmarkCount int64     `gorm:"default:0" json:"bookmark_count"`
	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	User    User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Post    *Post     `gorm:"foreignKey:PostID" json:"post"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"parent"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies"`
}
