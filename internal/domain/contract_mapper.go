package domain

import "job4j.ru/sharetrip-contract/internal/repository"

func toRepositoryCanCreateTripRequest(cliendId string) repository.ServiceCompanyRequest {
	return repository.ServiceCompanyRequest{
		CompanyId:   cliendId,
		ServiceCode: string(TripCreate),
	}
}
