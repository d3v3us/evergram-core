package validation

import (
	"fmt"
	"strings"

	"github.com/deveusss/evergram-core/common"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidationMiddleware is a middleware function to validate requests
func ValidationMiddleware(c *fiber.Ctx) error {
	var req interface{}

	// Parsing JSON request body
	if err := c.BodyParser(&req); err != nil {
		// If parsing JSON fails, try parsing as form data
		form := new(fiber.Map)
		if err := c.BodyParser(form); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request format"})
		}
		req = *form
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation error"})
		}

		// Handle multiple validation errors
		var errMsg strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			errMsg.WriteString(fmt.Sprintf("Field: %s failed validation for tag: %s\n", err.Field(), err.Tag()))
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: errMsg.String()})
	}

	// If validation succeeds, continue with the next handler
	return c.Next()
}

// formatValidationErrors formats the validation errors into a readable string
func formatValidationErrors(err error) string {
	var errorMsg string
	for _, err := range err.(validator.ValidationErrors) {
		errorMsg += fmt.Sprintf("Field: %s failed validation for tag: %s\n", err.Field(), err.Tag())
	}
	return errorMsg
}
