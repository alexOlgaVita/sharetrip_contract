package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/sharetrip-contract/internal/http/dto"
	"time"
)

type ContractRequest dto.ContractRequest

type CreateContractRequest dto.CreateContractRequest

type CreateContractResponse dto.CreateContractResponse

func (s *Server) CreateContract(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req CreateContractRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if err := createValidate(&req); err != nil {
		return err
	}

	resp, err := s.Service.CreateContract(ctx, dto.CreateContractRequest{
		CompanyId: req.CompanyId,
		StartsAt:  req.StartsAt,
		EndsAt:    req.EndsAt,
	})
	if err != nil {
		log.Errorw("s.ContractService.CreateTrip", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func createValidate(req *CreateContractRequest) error {
	if req.CompanyId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "companyId is required")
	}

	if req.StartsAt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "StartsAt is required")
	}
	now := time.Now()
	templateDate := "2006-01-02 15:04:05"
	startsAtTime, err := time.Parse(templateDate, req.StartsAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "startsAt: date format error")
	}
	if now.After(startsAtTime) {
		return fiber.NewError(fiber.StatusBadRequest, "startsAt is expired date")
	}

	if req.EndsAt != "" {
		endsAtTime, err := time.Parse(templateDate, req.EndsAt)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "endsAt: date format error")
		}
		if now.After(endsAtTime) {
			return fiber.NewError(fiber.StatusBadRequest, "endsAt is expired date")
		}
		if startsAtTime.After(endsAtTime) {
			return fiber.NewError(fiber.StatusBadRequest, "endsAt is before startsAt")
		}
	}

	return nil
}
