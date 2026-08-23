package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *Server) GetContract(c *fiber.Ctx) error {
	contractId := c.Params("contractId")

	if contractId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "contractId is required")
	}

	contract, err := s.Service.GetContract(c.Context(), contractId)

	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "trip is not found")
		}
		log.Errorw(
			"get trip failed",
			"error", err,
			"trip_id", contractId,
		)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := dto.Contract{
		ID:        contract.ID,
		CompanyId: contract.CompanyId,
		Status:    contract.Status,
		StartsAt:  contract.StartsAt,
		EndsAt:    contract.EndsAt,
		CreatedAt: contract.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(dto.GetContractResponse{Contract: res})
}
