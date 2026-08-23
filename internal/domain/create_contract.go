package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

func (u *ContractUsecase) CreateContract(
	ctx context.Context,
	tx pgx.Tx,
	req dto.CreateContractRequest,
) (*dto.Contract, error) {
	id := uuid.NewString()
	contract, err := u.ContractRepo.Create(ctx, dto.Contract{
		ID:        id,
		CompanyId: req.CompanyId,
		Status:    dto.ContractStatusDraft,
		StartsAt:  req.StartsAt,
		EndsAt:    req.EndsAt,
	})
	if err != nil {
		return nil, fmt.Errorf("repoContract.Create: %w", err)
	}
	return contract, nil
}
