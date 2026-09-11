package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"job4j.ru/sharetrip-contract/configs"
	"job4j.ru/sharetrip-contract/internal/domain"
	api "job4j.ru/sharetrip-contract/internal/http"
	"job4j.ru/sharetrip-contract/internal/repository"
	"job4j.ru/sharetrip-contract/internal/service"
	"log"
)

func main() {
	ctx := context.Background()

	cfg := repository.Config{
		Host:     configs.Env("DB_HOST", "localhost"),
		Port:     configs.EnvInt("DB_PORT", 6545),
		User:     configs.Env("DB_USER", "postgres"),
		Password: configs.Env("DB_PASSWORD", "password"),
		DBName:   configs.Env("DB_NAME", "sharetripContracts"),
		SSLMode:  configs.Env("DB_SSLMODE", "disable"),
	}

	pool, err := repository.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepoPg(pool)

	server := api.NewServer(repo)
	server.Service = &service.ContractService{
		Pool: pool,
		ContractUsecase: &domain.ContractUsecase{
			ContractRepo: repo,
		},
	}
	app := fiber.New()
	server.Route(app.Group("/"))

	err = app.Listen(":8082")
	if err != nil {
		log.Fatal(err)
	}

}
