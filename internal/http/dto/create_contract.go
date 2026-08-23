package dto

type CreateContractRequest struct {
	CompanyId string `json:"companyId"`
	StartsAt  string `json:"startsAt"`
	EndsAt    string `json:"endsAt"`
}

type CreateContractResponse struct {
	Contract ContractRequest `json:"contract"`
}
