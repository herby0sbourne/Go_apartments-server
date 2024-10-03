package main

import (
	"appartments-server/database"
	"appartments-server/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {
	godotenv.Load()
	database.InitializeDB()

	app := iris.Default()
	app.Validator = validator.New()

	location := app.Party("/api/location")
	{
		location.Get("/autocomplete", routes.Autocomplete)
		location.Get("/search", routes.Search)
	}

	user := app.Party("/api/user")
	{
		user.Post("/create-user", routes.Register)
		user.Post("/login-user", routes.Login)
		user.Post("/facebook-OAuth", routes.FacebookBookLoginOrSignUp)
	}

	app.Listen(":4000")
}
