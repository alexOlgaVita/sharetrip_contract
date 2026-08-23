package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

type ContractService struct {
	Pool            *pgxpool.Pool
	ContractUsecase *domain.ContractUsecase
}

func NewContractService(
	Pool *pgxpool.Pool,
	ContractUsecase *domain.ContractUsecase,
) *ContractService {
	return &ContractService{
		Pool:            Pool,
		ContractUsecase: ContractUsecase,
	}
}

func (s *ContractService) CreateContract(
	ctx context.Context,
	req dto.CreateContractRequest,
) (*dto.Contract, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Contract, error) {

		resp, err := s.ContractUsecase.CreateContract(ctx, tx, dto.CreateContractRequest{
			CompanyId: req.CompanyId,
			StartsAt:  req.StartsAt,
			EndsAt:    req.EndsAt,
		})
		if err != nil {
			return nil, fmt.Errorf("usecase.CreateContract: %w", err)
		}

		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
