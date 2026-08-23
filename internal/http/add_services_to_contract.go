package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *Server) AddServicesToContract(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req dto.ServicesContractRequest

	contractId := c.Params("contractId")
	if contractId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "contractId is required")
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if err := addValidate(&req); err != nil {
		return err
	}

	resp, err := s.Service.AddServicesToContract(ctx, contractId, req)
	if err != nil {
		if errors.Is(err, domain.ErrContractNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "contract is not found")
		}
		log.Errorw("s.ContractService.AddServicesToContract", err)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func addValidate(req *dto.ServicesContractRequest) error {
	if req.Services == nil {
		return fiber.NewError(fiber.StatusBadRequest, "services is required")
	}

	seen := make(map[string]struct{}, len(req.Services))

	for i := range req.Services {
		code := strings.TrimSpace(req.Services[i].ServiceCode)
		if code == "" {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("services[%d].serviceCode is required", i))
		}
		if len(code) > 64 {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("services[%d].serviceCode: max length is 64", i))
		}

		if _, exists := seen[code]; exists {
			return fiber.NewError(fiber.StatusConflict,
				fmt.Sprintf("duplicate serviceCode: %s", code))
		}
		seen[code] = struct{}{}
		req.Services[i].ServiceCode = code
	}

	return nil
}
