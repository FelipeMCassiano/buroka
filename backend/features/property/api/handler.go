package api

import (
	"fmt"

	"github.com/FelipeMCassiano/buroka/features/property/application"
	"github.com/FelipeMCassiano/buroka/features/property/domain"
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

	propertyCreated, err := h.service.RegisterProperty(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(propertyCreated)
}

func (h *PropertyHandler) GetProperty(c *fiber.Ctx) error {
	propertyName := c.Params("name")
	propertyCode := c.Params("code")

	property, err := h.service.GetProperty(propertyName, propertyCode)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(property)
}

func (h *PropertyHandler) SearchProperty(c *fiber.Ctx) error {
	propertyType := c.Query("type")
	forRent := c.QueryBool("forRent")
	forSale := c.QueryBool("forSale")
	neigborhood := c.Query("neigborhood")
	salePrice := c.QueryInt("price")
	rentAmount := c.QueryInt("rent")
	size := c.QueryFloat("size")

	searchFilter := domain.SearchFilter{
		PropertyType: propertyType,
		ForRent:      forRent,
		ForSale:      forSale,
		Neighborhood: neigborhood,
		RentAmount:   rentAmount,
		SalePrice:    salePrice,
		Size:         size,
	}

	properties, err := h.service.SearchProperty(searchFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(properties)
}
