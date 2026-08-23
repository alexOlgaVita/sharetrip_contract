package dto

type UpdateContractRequest struct {
	CompanyId string `json:"companyId"`
	StartsAt  string `json:"startsAt"`
	EndsAt    string `json:"endsAt"`

	ContractID string
	Status     string
}
