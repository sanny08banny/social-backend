package models


type UserDTO struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	ProfileName string `json:"profile_name"`
	ProfilePic  string `json:"profile_pic"`
}

type PostDTO struct {
	PostID       uint     `json:"post_id"`
	Content      string   `json:"content"`
	User         UserDTO  `json:"user"`
	DateCreated  string   `json:"date_created"`
	LastUpdated  string   `json:"last_updated"`
	CommentCount int64    `json:"comment_count"`
	LikeCount    int64    `json:"like_count"`
}