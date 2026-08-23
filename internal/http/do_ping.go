package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Server) DoPing(c *fiber.Ctx) error {
	err := s.Repository.DoPing(c.Context())
	if err != nil {
		log.Errorw("s.Repository.DoPing", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON("Все ок!")
}
