package dto

type ChangeContractsStatusModelRequest struct {
	Status string
}

type ChangeContractsStatusModelResponse struct {
	ID        string
	CompanyId string
	Status    string
	StartsAt  string
	EndsAt    string
	CreatedAt string
}
