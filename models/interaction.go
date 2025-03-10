package models

import "time"

type PostView struct {
	ViewID    uint      `gorm:"primaryKey" json:"view_id"`
	PostID    uint      `gorm:"index" json:"post_id"`
	UserID    uint      `json:"user_id"` // Nullable for guests
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Model for reposts
type Repost struct {
	RepostID       uint      `gorm:"primaryKey" json:"repost_id"`
	OriginalPostID uint      `json:"original_post_id"`
	UserID         uint      `json:"user_id"`
	RepostedAt     time.Time `gorm:"autoCreateTime" json:"reposted_at"`
}