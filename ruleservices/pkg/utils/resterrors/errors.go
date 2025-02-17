// Package resterrors provides functions to return different RESTful error responses.
package resterrors

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a RESTful error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse returns a new ErrorResponse instance.
func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{Code: code, Message: message}
}

// Send sends the ErrorResponse as a JSON response.
func (e ErrorResponse) Send(c *fiber.Ctx) error {
	return c.Status(e.Code).JSON(e)
}

// SendError sends a generic error as a JSON response.
func SendError(c *fiber.Ctx, code int, message string) error {
	return NewErrorResponse(code, message).Send(c)
}

// SendNotFoundError sends a 404 Not Found error as a JSON response.
func SendNotFoundError(c *fiber.Ctx, message string) error {
	return SendError(c, http.StatusNotFound, message)
}

// SendUnauthorizedError sends a 401 Unauthorized error as a JSON response.
func SendUnauthorizedError(c *fiber.Ctx) error {
	return SendError(c, http.StatusUnauthorized, "Unauthorized access")
}

// SendForbiddenError sends a 403 Forbidden error as a JSON response.
func SendForbiddenError(c *fiber.Ctx) error {
	return SendError(c, http.StatusForbidden, "Forbidden access")
}

// SendForbiddenError sends a 403 Forbidden error as a JSON response.
func SendBadRequestError(c *fiber.Ctx) error {
	return SendError(c, http.StatusBadRequest, "Invalid data, unable to parse")
}

// SendInternalServerError sends a 500 Internal Server Error as a JSON response.
func SendInternalServerError(c *fiber.Ctx) error {
	return SendError(c, http.StatusInternalServerError, "Internal server error")
}

// SendValidationError sends a 400 Bad Request error as a JSON response.
func SendValidationError(c *fiber.Ctx, message string) error {
	return SendError(c, http.StatusBadRequest, message)
}

func SendConflictError(c *fiber.Ctx, message string) error {
	return SendError(c, http.StatusConflict, message)
}

// SendCreated sends a 201 Created response.
func SendCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(data)
}

// SendOK sends a 200 OK response with data.
func SendOK(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(data)
}

// SendDeleted sends a 204 No Content response.
func SendNoContent(c *fiber.Ctx) error {
	return c.Status(http.StatusNoContent).SendString("")
}
