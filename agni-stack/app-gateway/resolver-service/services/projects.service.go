package services

import (
	"app-gateway/graph/model"
	graphmodel "app-gateway/graph/model"
	repository "app-gateway/repository"
	dbmodel "app-gateway/repository/model"
	"strconv"

	"context"
)

func CreateProject(ctx context.Context, input model.ProjectInput) (*graphmodel.Project, error) {

	project, err := repository.GetRepositoryManager().ProjectRepo.CreateProject(ctx, input)
	if err != nil {
		return nil, err
	}

	graphproj := mapProject(project)
	return graphproj, nil
}

func GetProjects(ctx context.Context) ([]*graphmodel.Project, error) {
	projects, err := repository.GetRepositoryManager().ProjectRepo.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return multiProjectMapper(projects), nil
}

func GetProject(ctx context.Context, project_id string) (*graphmodel.Project, error) {
	id, _ := strconv.ParseInt(project_id, 10, 64)
	project, err := repository.GetRepositoryManager().ProjectRepo.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	graphproj := mapProject(project)
	return graphproj, nil
}

func multiProjectMapper(projects []*dbmodel.Project) []*graphmodel.Project {

	graphProjects := make([]*graphmodel.Project, len(projects)) // Pre-allocate slice
	for i, dbProject := range projects {
		graphProjects[i] = mapProject(dbProject)
	}
	return graphProjects
}

func mapProject(project *dbmodel.Project) *graphmodel.Project {

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
		graphproj.Apps = getApps(project.Apps)
	}

	return graphproj
}

func getApps(apps []dbmodel.App) []*graphmodel.BasicApp {

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
