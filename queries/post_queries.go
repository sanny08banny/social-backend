package queries

const (
	GetPostsQuery = `
		SELECT post_id, user_id, content, date_created, last_updated, view_count, repost_count, comment_count, like_count, bookmark_count
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

	IncrementViewCountQuery = `
		UPDATE posts
		SET view_count = view_count + 1
		WHERE post_id = $1
	`

	IncrementRepostCountQuery = `
		UPDATE posts
		SET repost_count = repost_count + 1
		WHERE post_id = $1
	`
)
