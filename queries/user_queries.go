package queries

const (
	GetUsersQuery  = "SELECT id, name, email FROM users"
	CreateUserQuery = "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id"
)