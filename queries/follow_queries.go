package queries 

const (
	// Follow Queries
	CreateFollowQuery = `
		INSERT INTO follow (owner_id, user_id)
		VALUES ($1, $2)
		RETURNING follow_id, followed_at
	`

	DeleteFollowQuery = `
		DELETE FROM follow
		WHERE owner_id = $1 AND user_id = $2
	`
)