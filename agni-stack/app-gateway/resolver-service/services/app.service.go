package services

import (
	graphmodel "app-gateway/graph/model"
	repository "app-gateway/repository"
	dbmodel "app-gateway/repository/model"
	"app-gateway/utils"

	"context"
	"strconv"
)

/* Need to use transcation for the app creation as
this has a external request dependencies with
it. */

func CreateApp(ctx context.Context, input graphmodel.AppInput) (*graphmodel.App, error) {
	app, err := repository.GetRepositoryManager().AppRepo.CreateApp(ctx, input)
	if err != nil {
		return nil, err
	}

	err = kafkaMessage(*app)
	if err != nil {
		return nil, err
	}

	graphApp := mapApp(app)
	return graphApp, nil
}

func kafkaMessage(app dbmodel.App) error {
	message := map[string]string{
		"nodeApp":   app.Name,
		"nodeImage": app.Image,
	}

	producer := utils.NewProducer()
	err := producer.PushMessage(message)

	if err != nil {
		return err
	}

	return nil
}

func GetApps(ctx context.Context) ([]*graphmodel.App, error) {
	apps, err := repository.GetRepositoryManager().AppRepo.GetApps(ctx)
	if err != nil {
		return nil, err
	}

	graphApp := multiAppMapper(apps)
	return graphApp, nil
}

func GetApp(ctx context.Context, app_id string) (*graphmodel.App, error) {
	id, _ := strconv.ParseInt(app_id, 10, 64)
	app, err := repository.GetRepositoryManager().AppRepo.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}

	graphApp := mapApp(app)
	return graphApp, nil
}

func multiAppMapper(app []*dbmodel.App) []*graphmodel.App {

	graphApps := make([]*graphmodel.App, len(app)) // Pre-allocate slice
	for i, dbApp := range app {
		graphApps[i] = mapApp(dbApp)
	}
	return graphApps
}
func mapApp(app *dbmodel.App) *graphmodel.App {

	graphapp := &graphmodel.App{
		ID:          strconv.FormatInt(app.ID, 10),
		Name:        app.Name,
		Description: app.Description,
		Image:       &app.Image,
	}

	if app.Owner.ID != 0 {
		graphapp.Owner = &graphmodel.BasicUser{
			ID:    strconv.FormatInt(app.Owner.ID, 10),
			Name:  app.Owner.Name,
			Email: app.Owner.Email,
		}
	}

	if app.ProjectID != 0 {
		graphapp.Project = &graphmodel.BasicProject{
			ID:          strconv.FormatInt(app.ProjectID, 10),
			Name:        app.Project.Name,
			Description: app.Project.Description,
		}
	}

	return graphapp
}
