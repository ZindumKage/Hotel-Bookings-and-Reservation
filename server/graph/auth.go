package graph

import (
	"context"
	"fmt"
)

type AuthUser struct {
	ID   uint
	Role string
}
type contextKey string

const userContextKey = contextKey("user")

func GetUserFromCtx(ctx context.Context) (*AuthUser, error) {
	user, ok := ctx.Value(userContextKey).(*AuthUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}
