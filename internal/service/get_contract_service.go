package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *ContractService) GetContract(
	ctx context.Context,
	contractId string,
) (*dto.Contract, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Contract, error) {
		resp, err := s.ContractUsecase.GetContract(ctx, tx, contractId)
		if err != nil {
			return nil, fmt.Errorf("usecase.GetContract: %w", err)
		}
		return resp, nil

	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
