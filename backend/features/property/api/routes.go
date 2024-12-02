package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(pH *PropertyHandler, app *fiber.App) {
	app.Post("/properties/register-property", pH.RegisterProperty)
	app.Get("/properties/:name", pH.GetProperty)
}