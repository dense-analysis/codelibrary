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

	app.Post("/api/auth/login", routes.LoginHandler(db))
	app.Post("/api/auth/logout", routes.LogoutHandler(db))
	app.Post("/api/auth/register", routes.RegisterHandler(db))
	app.Get("/api/code", routes.ListCodeSamplesHandler(db))
	app.Post("/api/code", routes.CreateCodeSampleHandler(db))
	app.Get("/api/code/:id", routes.GetCodeSampleHandler(db))
	app.Put("/api/code/:id", routes.UpdateCodeSampleHandler(db))
	app.Delete("/api/code/:id", routes.DeleteCodeSampleHandler(db))
	app.Get("/api/docs/*", swagger.HandlerDefault)

	port := os.Getenv("API_PORT")

	if len(port) == 0 {
		port = "7000"
	}

	app.Listen(":" + port)
}
