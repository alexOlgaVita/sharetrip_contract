package domain

import (
	"context"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (u *ContractUsecase) ChangeContractsStatus(
	ctx context.Context,
	tx pgx.Tx,
	req dto.UpdateContractRequest,
) (*dto.Contract, error) {
	contract, err := u.ContractRepo.GetForUpdateByID(ctx, tx, req.ContractID)
	if err != nil {
		return nil, err
	}

	if req.Status == dto.ContractStatusActive {
		if contract.Status == dto.ContractStatusActive {
			return contract, ErrStatusIsPublishedAlready
		}
		if contract.Status != dto.ContractStatusDraft {
			return nil, ErrNotAllowedCurrentStatusToPublish
		}
	}
	err = u.ContractRepo.UpdateStatus(ctx, tx, contract.ID, dto.ContractStatusActive)
	if err != nil {
		return nil, err
	}
	contract.Status = dto.ContractStatusActive

	return contract, nil
}
