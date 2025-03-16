package services

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func GetFields(ctx context.Context) []string {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	var dataField []string

	for _, field := range fields {
		dataField = append(dataField, field.Name)
	}
	return dataField
}
