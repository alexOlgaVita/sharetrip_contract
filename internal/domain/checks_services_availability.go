package domain

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (u *ContractUsecase) ChecksServicesAvailability(
	ctx context.Context,
	tx pgx.Tx,
	req dto.ServiceCompanyRequest,
) (*dto.AvailbaleServicesCompanyResponse, error) {
	availbaleServicesCompany, err := u.ContractRepo.GetAvailabilityCompaniesService(ctx, tx, req)
	if err != nil {
		log.Errorw("s.Repository.GetAvailabilityCompaniesService", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return availbaleServicesCompany, nil
}
