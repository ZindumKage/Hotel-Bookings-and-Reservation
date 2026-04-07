package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/graph/model"
)

type contextKey string

const UserCtxKey = contextKey("user")

func roleRank(role model.Role) int {
	switch role {
	case model.RoleSuperAdmin:
		return 3
	case model.RoleAdmin:
		return 2
	case model.RoleUser:
		return 1
	default:
		return 0
	}
}

func HasRole(
	ctx context.Context,
	obj interface{},
	next graphql.Resolver,
	role model.Role,
) (interface{}, error) {

	user, ok := ctx.Value(UserCtxKey).(*model.User)
	if !ok {
		return nil, fmt.Errorf("unauthenticated")
	}

	if roleRank(user.Role) < roleRank(role) {
		return nil, fmt.Errorf("forbidden")
	}

	return next(ctx)
}