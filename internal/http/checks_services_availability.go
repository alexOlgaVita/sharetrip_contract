package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *Server) ChecksServicesAvailability(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req dto.ServiceCompanyRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if err := checkValidate(&req); err != nil {
		return err
	}

	resp, err := s.Service.ChecksServicesAvailability(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "contract is not found")
		}
		log.Errorw("s.ContractService.ChecksServicesAvailability", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func checkValidate(req *dto.ServiceCompanyRequest) error {
	if req.CompanyId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "companyId is required")
	}
	if _, err := uuid.Parse(req.CompanyId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "companyId must be a valid UUID")
	}
	if req.ServiceCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "serviceCode is required")
	}
	return nil
}
