package main

import (
	"log"

	"github.com/FelipeMCassiano/buroka/features/pkg/db"
	"github.com/FelipeMCassiano/buroka/features/property/api"
	"github.com/FelipeMCassiano/buroka/features/property/application"
	"github.com/FelipeMCassiano/buroka/features/property/infrastructure"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	repo := infrastructure.NewPropertyRepository(database)
	service := application.NewPropertyService(repo)
	handler := api.NewPropertyHandler(service)

	api.RegisterRoutes(handler, app)

	app.Listen(":3000")
}
