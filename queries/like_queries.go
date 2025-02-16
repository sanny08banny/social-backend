package queries

const (
	GetLikesQuery    = "SELECT * FROM likes"
    GetLikeByIDQuery = "SELECT * FROM likes WHERE like_id = $1"
    CreateLikeQuery  = `INSERT INTO likes (post_id, user_id) VALUES ($1, $2) RETURNING like_id`
    DeleteLikeQuery  = "DELETE FROM likes WHERE like_id = $1"
)