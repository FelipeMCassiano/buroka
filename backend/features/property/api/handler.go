package api

import (
	"fmt"

	"github.com/FelipeMCassiano/buroka/features/property/application"
	"github.com/gofiber/fiber/v2"
)

type PropertyHandler struct {
	service *application.PropertyService
}

func NewPropertyHandler(service *application.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: service}
}

// handlers here

func (h *PropertyHandler) RegisterProperty(c *fiber.Ctx) error {
	req := new(application.RegisterPropertyRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.service.RegisterProperty(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *PropertyHandler) GetProperty(c *fiber.Ctx) error {
	propertyName := c.Params("name")

	property, err := h.service.GetProperty(propertyName)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(property)
}
