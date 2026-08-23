package domain

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (u *ContractUsecase) GetContract(
	ctx context.Context,
	tx pgx.Tx,
	contractId string,
) (*dto.Contract, error) {
	contract, err := u.ContractRepo.GetByID(ctx, tx, contractId)
	if err != nil {
		log.Errorw("s.Repository.Get", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return contract, nil
}
