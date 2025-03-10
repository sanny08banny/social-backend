package queries

const (
	CreatePostViewQuery = `
		INSERT INTO post_views (post_id, user_id, ip_address)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, user_id, ip_address) DO NOTHING
	`

	CheckPostViewExistsQuery = `
		SELECT COUNT(*) FROM post_views WHERE post_id = $1 AND (user_id = $2 OR ip_address = $3)
	`
)
