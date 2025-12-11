package core

import "github.com/gofiber/fiber/v2"

type Module interface {
	Register(container *Container)
	MountRoutes(router fiber.Router)
}
