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

	route.Post(
		"contracts/service-availability-checks",
		s.ChecksServicesAvailability,
	)

	route.Put(
		"/contracts/:contractId/services",
		s.AddServicesToContract,
	)
}
