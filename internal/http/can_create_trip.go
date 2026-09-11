package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"job4j.ru/sharetrip-contract/internal/domain"
)

func (s *Server) CanCreateTrip(c *fiber.Ctx) error {
	ctx := c.UserContext()

	clientId := c.Params("client_id")
	if err := checkValidate(clientId); err != nil {
		return err
	}

	resp, err := s.Service.CanCreateTrip(ctx, clientId)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "contract is not found")
		}
		log.Errorw("s.ContractService.ChecksServicesAvailability", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func checkValidate(companyId string) error {
	if companyId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "companyId is required")
	}
	if _, err := uuid.Parse(companyId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "companyId must be a valid UUID")
	}
	return nil
}
