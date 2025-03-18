package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// type User struct {
// 	UserID       uint      `gorm:"primaryKey" json:"user_id"`
// 	Username     string    `gorm:"uniqueIndex" json:"username"`
// 	ProfileName  string    `json:"profile_name"`
// 	Email        string    `gorm:"uniqueIndex" json:"email"`
// 	Password     string    `json:"-"` // Exclude password from JSON responses
// 	Bio          string    `json:"bio"`
// 	PhoneNumber  string    `json:"phone_number"`
// 	ProfilePic   string    `json:"profile_pic"`
// 	OnlineStatus string    `json:"online_status"`
// 	DateCreated  time.Time `gorm:"autoCreateTime" json:"date_created"`
// 	LastUpdated  time.Time `gorm:"autoUpdateTime" json:"last_updated"`

// 	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts"`
// 	Comments []Comment `gorm:"foreignKey:UserID" json:"comments"`
// 	Likes    []Like    `gorm:"foreignKey:UserID" json:"likes"`
// }

type User struct {
	UserID       uint      `gorm:"primaryKey" json:"user_id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	ProfileName  string    `json:"profile_name"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Password     string    `gorm:"not null" json:"-"`
	Bio          string    `json:"bio"`
	PhoneNumber  string    `json:"phone_number"`
	ProfilePic   string    `json:"profile_pic"`
	OnlineStatus string    `gorm:"default:'offline'" json:"online_status"`
	DateCreated  time.Time `gorm:"autoCreateTime" json:"date_created"`
	LastUpdated  time.Time `gorm:"autoUpdateTime" json:"last_updated"`
	Followers    uint      `gorm:"default:0" json:"followers"`
	Following    uint      `gorm:"default:0" json:"following"`
}


type NewUser struct {
	UserID       uint      `gorm:"primaryKey" json:"user_id"`
	Username     string    `gorm:"uniqueIndex" json:"username"`
	ProfileName  string    `json:"profile_name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Password     string    `json:"password"` // Exclude password from JSON responses
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

// Hash password before saving to the database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword compares the hashed password with the provided password
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
