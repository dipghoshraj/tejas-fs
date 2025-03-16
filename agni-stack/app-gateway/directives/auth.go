package directives

import (
	"app-gateway/utils"
	"context"
	"errors"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

var jwtSecret = []byte("secret")

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

	httpRequest := graphql.GetOperationContext(ctx)
	if httpRequest == nil {
		return nil, errors.New("unable to retrieve request context")
	}

	// Get the token from the context
	authHeader := httpRequest.Headers.Get("Authorization")

	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, errors.New("invalid authorization format")
	}

	// Verify the token
	userid, err := utils.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Inject user ID into context for resolvers
	ctx = context.WithValue(ctx, "user_id", userid)

	// Call the next resolver
	return next(ctx)
}
