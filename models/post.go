package models

import "time"

// type Post struct {
// 	PostID      uint      `gorm:"primaryKey" json:"post_id"`
// 	UserID      uint      `gorm:"index" json:"user_id"`
// 	Content     string    `json:"content"`
// 	// ImageURL    string    `json:"image_url"`
// 	DateCreated time.Time `gorm:"autoCreateTime" json:"date_created"`
// 	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`

// 	User     User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
// 	Comments []Comment `gorm:"foreignKey:PostID" json:"comments"`
// 	Likes    []Like    `gorm:"foreignKey:PostID" json:"likes"`
// 	BookMarks    []BookMark    `gorm:"foreignKey:PostID" json:"bookmarks"`

// }
type Post struct {
	PostID        uint      `gorm:"primaryKey" json:"post_id"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Content       string    `json:"content"`
	DateCreated   time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated   time.Time `gorm:"autoUpdateTime" json:"last_updated"`
	ViewCount     int64     `gorm:"default:0" json:"view_count"`
	RepostCount   int64     `gorm:"default:0" json:"repost_count"`
	CommentCount   int64     `gorm:"default:0" json:"comment_count"`
	LikeCount     int64     `gorm:"default:0" json:"like_count"`
	BookmarkCount int64     `gorm:"default:0" json:"bookmark_count"`
	IsLiked       bool      `gorm:"-" json:"is_liked"`       // Dynamically computed
	IsBookmarked  bool      `gorm:"-" json:"is_bookmarked"`
}
