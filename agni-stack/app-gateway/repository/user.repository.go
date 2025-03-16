package repository

import (
	"app-gateway/repository/database"
	dbmodel "app-gateway/repository/model"
	"context"
	"fmt"
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user dbmodel.User) (*dbmodel.User, error)
	GetUser(ctx context.Context, id int64) (*dbmodel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*dbmodel.User, error)
}

func (r *userRepo) CreateUser(ctx context.Context, user dbmodel.User) (*dbmodel.User, error) {

	err := database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*dbmodel.User, error) {
	user := dbmodel.User{}
	query := database.DB.Where("users.id = ?", id)
	fields := GetFields(ctx)

	fmt.Printf("%v", fields)

	// Need to cap the preload of the apps and projects with a limit

	for _, preloadField := range []string{"projects"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField)).Preload("Projects.Apps")
		}
	}

	err := query.Omit("password").First(&user).Error
	if err != nil {
		fmt.Printf("err : %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*dbmodel.User, error) {
	user := dbmodel.User{}
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
