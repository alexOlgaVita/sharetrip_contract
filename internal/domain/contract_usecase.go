package domain

import "job4j.ru/sharetrip-contract/internal/repository"

type ContractUsecase struct {
	ContractRepo *repository.RepoPg
}
