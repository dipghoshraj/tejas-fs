package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"app-gateway/graph"
	"app-gateway/graph/model"
	"context"
	"fmt"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (string, error) {
	panic(fmt.Errorf("not implemented: Login - login"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
