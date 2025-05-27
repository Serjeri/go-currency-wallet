package models

type User struct {
	Name     string `form:"name"`
	Password string `form:"password"`
	Email    string `form:"email"`
}
// DB_HOST=localhost
// DB_PORT=5432
// DB_USER=postgres
// DB_PASSWORD=12051988
// DB_NAME=pasebin
