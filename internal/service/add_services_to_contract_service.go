package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *ContractService) AddServicesToContract(
	ctx context.Context,
	contractID string,
	req dto.ServicesContractRequest,
	// ) (*dto.ServicesContract, error) {
) (*[]dto.ContractService, error) {
	//	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.ServicesContract, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*[]dto.ContractService, error) {
		resp, err := s.ContractUsecase.AddServicesToContract(ctx, tx, contractID, req)
		if err != nil {
			return nil, fmt.Errorf("usecase.AddServicesToContract: %w", err)
		}
		return resp, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
