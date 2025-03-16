package repository

import (
	graphmodel "app-gateway/graph/model"
	"app-gateway/repository/database"
	dbmodel "app-gateway/repository/model"
	"context"
	"fmt"
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ProjectRepo interface {
	CreateProject(ctx context.Context, input graphmodel.ProjectInput) (*dbmodel.Project, error)
	GetProjects(ctx context.Context) ([]*dbmodel.Project, error)
	GetProject(ctx context.Context, id int64) (*dbmodel.Project, error)
}

func (pr *projRepo) CreateProject(ctx context.Context, input graphmodel.ProjectInput) (*dbmodel.Project, error) {

	userID, ok := ctx.Value("user_id").(float64)
	fields := GetFields(ctx)

	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	project := &dbmodel.Project{
		Name:        input.Name,
		OwnerID:     int64(userID),
		Description: *input.Description,
	}

	err := database.DB.Create(&project).Error
	if err != nil {
		return project, err
	}

	if slices.Contains(fields, "owner") {
		database.DB.Preload("Owner").First(&project, project.ID)
	}
	return project, nil
}

func (pr *projRepo) GetProjects(ctx context.Context) ([]*dbmodel.Project, error) {
	var projects []*dbmodel.Project

	userID, ok := ctx.Value("user_id").(float64)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	fields := GetFields(ctx)
	query := database.DB.Where("projects.owner_id = ?", userID)

	for _, preloadField := range []string{"owner", "apps"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField))
		}
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (pr *projRepo) GetProject(ctx context.Context, id int64) (*dbmodel.Project, error) {

	var dbproject dbmodel.Project

	userID, ok := ctx.Value("user_id").(float64)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	fields := GetFields(ctx)
	query := database.DB.Where("projects.owner_id = ? AND projects.ID = ?", userID, id)

	for _, preloadField := range []string{"owner", "apps"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField))
		}
	}

	if err := query.First(&dbproject).Error; err != nil {
		return nil, err
	}

	return &dbproject, nil
}
