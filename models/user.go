package models

import "time"

type User struct {
	UserID       uint      `gorm:"primaryKey" json:"user_id"`
	Username     string    `gorm:"uniqueIndex" json:"username"`
	ProfileName  string    `json:"profile_name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Bio          string    `json:"bio"`
	PhoneNumber  string    `json:"phone_number"`
	ProfilePic   string    `json:"profile_pic"`
	OnlineStatus string    `json:"online_status"`
	DateCreated  time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated  time.Time `gorm:"autoUpdateTime" json:"last_updated"`

	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:UserID" json:"likes"`
}
