package api

import (
	"job4j.ru/sharetrip-contract/internal/repository"
	"job4j.ru/sharetrip-contract/internal/service"
)

type Server struct {
	Repository *repository.RepoPg
	Service    *service.ContractService
}

func NewServer(repo *repository.RepoPg,
) *Server {
	s := &Server{
		Repository: repo,
	}

	return s
}
