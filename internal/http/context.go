package http

import (
	"context"
)

type contextKey string

const userKey contextKey = "user"

type User struct {
	ID int
	Name string
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}