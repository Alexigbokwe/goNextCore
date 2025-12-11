package core

import "github.com/gofiber/fiber/v2"

// Middleware interface for defining standardized middleware
type Middleware interface {
	Use() fiber.Handler
}

// HandlerMiddleware wraps a standard Fiber handler into a Middleware
type HandlerMiddleware struct {
	Handler fiber.Handler
}

func (h HandlerMiddleware) Use() fiber.Handler {
	return h.Handler
}

// Combine combines multiple middlewares into a fiber handler slice
func Combine(middlewares ...Middleware) []fiber.Handler {
	handlers := make([]fiber.Handler, len(middlewares))
	for i, m := range middlewares {
		handlers[i] = m.Use()
	}
	return handlers
}
