package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

type ChangeContractsStatusModelRequest dto.ChangeContractsStatusModelRequest

func (s *Server) ChangeContractsStatus(c *fiber.Ctx) error {
	var req ChangeContractsStatusModelRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	contractId := c.Params("contractId")
	if contractId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "contractId is required")
	}

	if _, err := uuid.Parse(contractId); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "contractId hasn't uuid format")
	}

	var resp, err = s.Service.ChangeContractsStatus(c.UserContext(), dto.UpdateContractRequest{
		ContractID: contractId,
		Status:     req.Status,
	})

	if err != nil {

		if errors.Is(err, domain.ErrContractNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "contract is not found")
		}
		if errors.Is(err, domain.ErrNotAllowedCurrentStatusToPublish) {
			return fiber.NewError(fiber.StatusConflict, "current status is not allowed for publish")
		}
		if errors.Is(err, domain.ErrStatusIsPublishedAlready) {
			return fiber.NewError(fiber.StatusNoContent, "contract's status is published already")
		}

		log.Errorw("s.service.MoveContractDraftToActive", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
