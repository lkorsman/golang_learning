package auth

import (
	"context"
)

type User struct {
	ID 		int 	`json:"id"`
	Email	string	`json:"email"`
	Password string `json:"-"`
}

type contextKey string

const userKey contextKey = "user"

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok 
}

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}