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

type Follow struct {
	FollowID   uint      `gorm:"primaryKey" json:"follow_id"`
	OwnerID    uint      `json:"owner_id"`
	UserID     uint      `json:"user_id"`
	FollowedAt time.Time `gorm:"autoCreateTime" json:"followed_at"`

	OwnerUser User `gorm:"foreignKey:OwnerID;references:UserID;constraint:OnDelete:CASCADE;" json:"owner_user"`
	User      User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;" json:"user"`
}
