package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready/", s.DoPing)

	route.Post(
		"/contracts/",
		s.CreateContract,
	)

	route.Put(
		"/contracts/:contractId/status-changes",
		s.ChangeContractsStatus,
	)

	route.Get(
		"/contracts/:contractId",
		s.GetContract,
	)

	route.Get(
		"contracts/can_create_trip/:client_id",
		s.CanCreateTrip,
	)

	route.Put(
		"/contracts/:contractId/services",
		s.AddServicesToContract,
	)
}
