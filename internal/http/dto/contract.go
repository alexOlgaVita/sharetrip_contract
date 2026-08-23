package dto

type ContractRequest struct {
	ID        string `json:"id"`
	CompanyId string `json:"companyId"`
	Status    string `json:"status"`
	StartsAt  string `json:"StartsAt"`
	EndsAt    string `json:"EndsAt"`
}

type Contract struct {
	ID        string
	CompanyId string
	Status    string
	StartsAt  string
	EndsAt    string
	CreatedAt string
}

type GetContractResponse struct {
	Contract Contract `json:"contract"`
}
