package queries

const (
	GetPostsQuery = `
		SELECT post_id, user_id, content, date_created, last_updated
		FROM posts
	`

	CreatePostQuery = `
		INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING post_id, date_created, last_updated
	`

	UpdatePostQuery = `
		UPDATE posts
		SET content = $1, last_updated = NOW()
		WHERE post_id = $2
	`

	DeletePostQuery = `
		DELETE FROM posts
		WHERE post_id = $1
	`
)