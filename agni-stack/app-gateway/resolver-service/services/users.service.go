package services

import (
	"app-gateway/graph/model"
	graphmodel "app-gateway/graph/model"
	repository "app-gateway/repository"
	dbmodel "app-gateway/repository/model"
	"app-gateway/utils"
	"context"
	"fmt"
	"strconv"
)

func CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	hashpassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := dbmodel.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashpassword,
	}

	userData, err := repository.GetRepositoryManager().UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// userData :=

	return &model.User{
		ID:    strconv.FormatInt(userData.ID, 10),
		Name:  userData.Name,
		Email: userData.Email,
	}, nil
}

func LoginUser(ctx context.Context, input model.LoginInput) (*model.AuthResponse, error) {

	user, err := repository.GetRepositoryManager().UserRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	valid_password := utils.CheckPassword(user.Password, input.Password)
	if !valid_password {
		return nil, fmt.Errorf("Invalid password")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		ID:    strconv.FormatInt(user.ID, 10),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := repository.GetRepositoryManager().UserRepo.GetUser(ctx, id)

	if err != nil {
		return nil, err
	}

	user_object := &model.User{
		ID:    strconv.FormatInt(user.ID, 10),
		Name:  user.Name,
		Email: user.Email,
	}

	/* TECH-DEBT : this section is very inefficient this need to
	optimised with better model design so we dont
	have to iterate, no feeling for doing it right now
	will do it after setup Kafka, for better implementation
	can look at project
	*/

	if len(user.Projects) != 0 {
		user_object.Projects = userProjects(user.Projects)
	}

	return user_object, nil
}

func userProjects(projects []dbmodel.Project) []*graphmodel.Project {

	graphProjects := make([]*graphmodel.Project, len(projects)) // Pre-allocate slice
	for i, dbProject := range projects {
		graphProjects[i] = structProject(dbProject)
	}
	return graphProjects
}

func structProject(project dbmodel.Project) *graphmodel.Project {

	graphproj := &graphmodel.Project{
		ID:          strconv.FormatInt(project.ID, 10),
		Name:        project.Name,
		Description: project.Description,
	}

	if project.Owner.ID != 0 {
		graphproj.Owner = &graphmodel.BasicUser{
			ID:    strconv.FormatInt(project.Owner.ID, 10),
			Name:  project.Owner.Name,
			Email: project.Owner.Email,
		}
	}

	if project.Apps != nil {
		graphproj.Apps = structApps(project.Apps)
	}

	return graphproj
}

func structApps(apps []dbmodel.App) []*graphmodel.BasicApp {

	basicApps := make([]*graphmodel.BasicApp, len(apps))

	for i, app := range apps {
		basicApps[i] = &graphmodel.BasicApp{
			ID:          strconv.FormatInt(app.ID, 10),
			Name:        app.Name,
			Image:       &app.Image,
			Description: &app.Description,
		}
	}

	return basicApps
}
