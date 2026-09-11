package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *ContractService) CanCreateTrip(
	ctx context.Context,
	clientId string,
) (*dto.AvailbaleServicesCompanyResponse, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.AvailbaleServicesCompanyResponse, error) {
		resp, err := s.ContractUsecase.CanCreateTrip(ctx, tx, clientId)
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
