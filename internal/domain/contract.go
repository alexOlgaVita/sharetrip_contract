package domain

type ContractService string

const (
	TripStarted ContractService = "trip.started"
	TripCreate  ContractService = "trip.create"
)
