package dto

type AvailbaleServicesCompanyResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

const (
	ReasonNoActiveContract     = "NO_ACTIVE_CONTRACT"
	ReasonContractSuspended    = "CONTRACT_SUSPENDED"
	ReasonContractTerminated   = "CONTRACT_TERMINATED"
	ReasonContractNotStarted   = "CONTRACT_NOT_STARTED"
	ReasonContractExpired      = "CONTRACT_EXPIRED"
	ReasonServiceNotInContract = "SERVICE_NOT_IN_CONTRACT"
	ReasonServiceDisabled      = "SERVICE_DISABLED"
)
