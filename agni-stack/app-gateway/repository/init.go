package repository

import (
	"context"
	"log"
	"sync"

	"github.com/99designs/gqlgen/graphql"
)

var instance *RepositoryManager
var once sync.Once

type userRepo struct{}
type projRepo struct{}
type appRepo struct{}

type RepositoryManager struct {
	UserRepo    UserRepo
	ProjectRepo ProjectRepo
	AppRepo     AppRepo
}

func NewUserRepository() UserRepo {
	return &userRepo{}
}

func NewProjectRepository() ProjectRepo {
	return &projRepo{}
}

func NewAppRepository() AppRepo {
	return &appRepo{}
}

/*
using singletone pattern for repository manager
to avide the multiple dependency injection of repository manager
to keep the single instance of repository manager

make code more readable and maintainable
*/

func InitRepositoryManager() {
	once.Do(func() {
		instance = &RepositoryManager{
			UserRepo:    NewUserRepository(),
			ProjectRepo: NewProjectRepository(),
			AppRepo:     NewAppRepository(),
		}
		log.Println("RepositoryManager initialized")
	})
}

func GetRepositoryManager() *RepositoryManager {
	if instance == nil {
		log.Fatal("RepositoryManager is not initialized")
	}

	return instance
}

func GetFields(ctx context.Context) []string {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	var dataField []string

	for _, field := range fields {
		dataField = append(dataField, field.Name)
	}
	return dataField
}
