package resolver

import (
	"app-gateway/database"
	"app-gateway/graph/model"
	dbmodel "app-gateway/resolver-service/model"
	"app-gateway/utils"
	"context"
	"fmt"
)

func CreateUser(ctx context.Context, input model.UserInput) (*dbmodel.User, error) {
	hashpassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := dbmodel.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashpassword,
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(ctx context.Context, input model.LoginInput) (string, error) {
	user := dbmodel.User{}
	err := database.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return "", err
	}

	valid_password := utils.CheckPassword(user.Password, input.Password)
	if !valid_password {
		return "", fmt.Errorf("Invalid password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUser(ctx context.Context, id int64) (*dbmodel.User, error) {
	user := dbmodel.User{}
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
