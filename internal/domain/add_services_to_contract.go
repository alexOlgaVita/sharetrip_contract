package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-contract/internal/http/dto"
	"job4j.ru/sharetrip-contract/internal/repository"
)

func (u *ContractUsecase) AddServicesToContract(
	ctx context.Context,
	tx pgx.Tx,
	contractID string,
	req dto.ServicesContractRequest,
	// ) (*dto.ServicesContract, error) {
) (*[]dto.ContractService, error) {
	_, err := u.ContractRepo.GetByID(ctx, tx, contractID)
	if err != nil {
		if errors.Is(err, repository.ErrContractNotFound) {
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("repoContract.GetByID: %w", err)
	}

	var services []dto.ContractService
	for _, item := range req.Services {
		service, err := u.ContractRepo.CreateContractService(ctx, dto.ContractService{
			ID:          uuid.NewString(),
			ContractId:  contractID,
			ServiceCode: item.ServiceCode,
			IsAvailable: item.Enabled,
		})
		if err != nil {
			return nil, fmt.Errorf("repoContract.CreateContractService: %w", err)
		}
		services = append(services, *service)
	}

	return &services, nil
}
