package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(pH *PropertyHandler, app *fiber.App) {
	propertiesRouter := app.Group("/properties")
	propertiesRouter.Post("/register-property", pH.RegisterProperty)
	propertiesRouter.Get("/:name/:code", pH.GetProperty)
	propertiesRouter.Get("/search", pH.SearchProperty)
}
