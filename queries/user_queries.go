package queries

const (
	GetUsersQuery = `
		SELECT user_id, username, profile_name, email, bio, phone_number, profile_pic, online_status, date_created, last_updated
		FROM users
	`

	CreateUserQuery = `
    INSERT INTO users (username, profile_name, email, password, bio, phone_number, profile_pic, online_status)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING user_id, date_created, last_updated
	`

	UpdateUserQuery = `
		UPDATE users
		SET username = $1, profile_name = $2, email = $3, bio = $4, phone_number = $5, profile_pic = $6, online_status = $7, last_updated = NOW()
		WHERE user_id = $8
	`

	DeleteUserQuery = `
		DELETE FROM users
		WHERE user_id = $1
	`
)