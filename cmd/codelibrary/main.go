package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/swagger"

	"github.com/dense-analysis/codelibrary/internal/api/database"
	"github.com/dense-analysis/codelibrary/internal/api/errorhandler"
	"github.com/dense-analysis/codelibrary/internal/api/routes"
	_ "github.com/dense-analysis/codelibrary/internal/docs"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.ErrorHandler,
	})
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))
	db, err := database.New(context.Background())

	if err != nil {
		// TODO: Exit more gracefully.
		panic(err)
	}

	app.Post("/auth/login", routes.LoginHandler(db))
	app.Post("/auth/logout", routes.LogoutHandler(db))
	app.Post("/auth/register", routes.RegisterHandler(db))
	app.Get("/code", routes.ListCodeSamplesHandler(db))
	app.Post("/code", routes.CreateCodeSampleHandler(db))
	app.Get("/code/:id", routes.GetCodeSampleHandler(db))
	app.Put("/code/:id", routes.UpdateCodeSampleHandler(db))
	app.Delete("/code/:id", routes.DeleteCodeSampleHandler(db))
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Listen(":8080")
}
