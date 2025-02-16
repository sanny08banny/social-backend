package queries

const (
	GetCommentsQuery    = "SELECT * FROM comments"
	GetCommentByIDQuery = "SELECT * FROM comments WHERE comment_id = $1"
	CreateCommentQuery  = `INSERT INTO comments (post_id, user_id, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING comment_id`
	UpdateCommentQuery  = `UPDATE comments SET content=$1, last_updated=NOW() WHERE comment_id=$2`
	DeleteCommentQuery  = "DELETE FROM comments WHERE comment_id = $1"
)
