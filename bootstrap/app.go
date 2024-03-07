package bootstrap

import "github.com/muhafs/go-restaurant-management/mongodb"

type Application struct {
	Env   *Env
	Mongo mongodb.Client
}

func App() Application {
	app := new(Application)
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
