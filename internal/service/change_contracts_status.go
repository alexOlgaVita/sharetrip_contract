package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/domain"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (s *ContractService) ChangeContractsStatus(
	ctx context.Context,
	req dto.UpdateContractRequest,
) (*dto.Contract, error) {
	res, err := tx(ctx, s.Pool, func(tx pgx.Tx) (*dto.Contract, error) {
		resp, err := s.ContractUsecase.ChangeContractsStatus(ctx, tx, dto.UpdateContractRequest{
			ContractID: req.ContractID,
			Status:     req.Status,
		})
		if err != nil {
			if errors.Is(err, domain.ErrContractNotFound) ||
				errors.Is(err, domain.ErrClientNotDriver) ||
				errors.Is(err, domain.ErrNotAllowedCurrentStatusToPublish) ||
				errors.Is(err, domain.ErrStatusIsPublishedAlready) {
				return nil, err
			}
			return nil, fmt.Errorf("usecase.ChangeContractsStatus: %w", err)
		}

		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
