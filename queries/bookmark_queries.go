package queries

const (
	GetBookmarksQuery    = "SELECT * FROM bookmarks"
	GetBookmarkByIDQuery = "SELECT * FROM bookmarks WHERE bookmark_id = $1"
	CreateBookmarkQuery  = `
		INSERT INTO bookmarks (user_id, post_id)
		VALUES ($1, $2)
		RETURNING bookmark_id
	`
	DeleteBookmarkQuery = `
		DELETE FROM bookmarks
		WHERE bookmark_id = $1
	`

	UpdateBookmarkCountQuery = `
		UPDATE posts
		SET bookmark_count = (
			SELECT COUNT(*) FROM bookmarks WHERE post_id = $1
		)
		WHERE post_id = $1
	`
)
