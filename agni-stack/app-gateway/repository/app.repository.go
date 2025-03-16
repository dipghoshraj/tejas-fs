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

type AppRepo interface {
	CreateApp(ctx context.Context, input graphmodel.AppInput) (*dbmodel.App, error)
	GetApp(ctx context.Context, id int64) (*dbmodel.App, error)
	GetApps(ctx context.Context) ([]*dbmodel.App, error)
}

func (ar *appRepo) CreateApp(ctx context.Context, input graphmodel.AppInput) (*dbmodel.App, error) {

	userID, ok := ctx.Value("user_id").(float64)
	fields := GetFields(ctx)

	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	app := &dbmodel.App{
		Name:        input.Name,
		Description: *input.Description,
		Image:       *input.Image,
		ProjectID:   int64(*input.ProjectID),
		OwnerID:     int64(userID),
	}

	err := database.DB.Create(&app).Error
	if err != nil {
		return nil, err
	}

	query := database.DB

	for _, preloadField := range []string{"owner", "project"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField))
		}
	}

	if err := query.Find(&app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func (ar *appRepo) GetApps(ctx context.Context) ([]*dbmodel.App, error) {
	var apps []*dbmodel.App
	fields := GetFields(ctx)

	userID, ok := ctx.Value("user_id").(float64)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	query := database.DB.Where("apps.owner_id = ?", userID)

	for _, preloadField := range []string{"owner", "project"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField))
		}
	}

	if err := query.Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil

}

func (ar *appRepo) GetApp(ctx context.Context, id int64) (*dbmodel.App, error) {

	var app *dbmodel.App
	fields := GetFields(ctx)

	userID, ok := ctx.Value("user_id").(float64)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	query := database.DB.Where("apps.owner_id = ? and apps.id = ?", userID, id)

	for _, preloadField := range []string{"owner", "project"} {
		if slices.Contains(fields, preloadField) {
			query = query.Preload(cases.Title(language.English).String(preloadField))
		}
	}

	if err := query.First(&app).Error; err != nil {
		return nil, err
	}

	return app, nil

}
