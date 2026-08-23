package dto

type ContractServiceItemRequest struct {
	ServiceCode string `json:"serviceCode"`
	Enabled     bool   `json:"enabled"`
}

type ServicesContractRequest struct {
	Services []ContractServiceItemRequest `json:"services"`
}

type ContractServiceItemResponse struct {
	ID          string `json:"id"`
	ContractId  string `json:"contractId"`
	ServiceCode string `json:"serviceCode"`
	IsAvailable bool   `json:"isAvailable"`
	CreatedAt   string `json:"createdAt"`
}

type ServicesContractResponse struct {
	Services []ContractServiceItemResponse `json:"services"`
}
