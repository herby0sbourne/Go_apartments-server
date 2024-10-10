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
		location.Get("/search", routes.Search)
		location.Get("/autocomplete", routes.Autocomplete)
	}

	user := app.Party("/api/user")
	{
		user.Post("/login-user", routes.Login)
		user.Post("/create-user", routes.Register)
		user.Post("/apple-login", routes.AppleLoginOrSignUp)
		user.Post("/google-login", routes.GoogleLoginOrSignUp)
		user.Post("/facebook-login", routes.FacebookBookLoginOrSignUp)
	}

	app.Listen(":4000")
}
