package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/sharetrip-contract/internal/http/dto"
)

type RepoPg struct {
	pool *pgxpool.Pool
}

func NewRepoPg(pool *pgxpool.Pool) *RepoPg {

	return &RepoPg{pool: pool}
}

func (r *RepoPg) Create(ctx context.Context, it dto.Contract) (*dto.Contract, error) {
	fmt.Println("id = "+it.ID+", company_id = "+it.CompanyId, ", status = "+it.Status, ", starts_at = "+it.StartsAt, ", ends_at = "+it.EndsAt)
	_, err := r.pool.Exec(
		ctx,
		`insert into contracts(id, company_id, status, starts_at, ends_at) values($1, $2, $3, $4, $5)`,
		it.ID, it.CompanyId, it.Status, it.StartsAt, it.EndsAt,
	)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Exec: %w", err)
	}

	return &it, nil
}

func (r *RepoPg) List(ctx context.Context) ([]dto.Contract, error) {
	rows, err := r.pool.Query(ctx, `select id, companyId from contracts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []dto.Contract
	for rows.Next() {
		var item dto.Contract
		if err := rows.Scan(&item.ID, &item.CompanyId, &item.Status, &item.StartsAt, &item.EndsAt); err != nil {
			return nil, err
		}
		contracts = append(contracts, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (r *RepoPg) Get(ctx context.Context, contractId string) (dto.Contract, error) {
	var it dto.Contract
	err := r.pool.QueryRow(
		ctx,
		`select id, company_id, status, 
       COALESCE(to_char(starts_at, 'MM-DD-YYYY HH24:MI'), ''), 
       COALESCE(to_char(ends_at, 'MM-DD-YYYY HH24:MI'), '') from contracts where id = $1`,
		contractId,
	).Scan(&it.ID, &it.CompanyId, &it.Status, &it.StartsAt, &it.EndsAt, &it.Status)

	return it, err
}

func (r *RepoPg) GetByID(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*dto.Contract, error) {
	contact := &dto.Contract{}

	err := tx.QueryRow(
		ctx,
		`select id, company_id, status, 
       COALESCE(to_char(starts_at, 'MM-DD-YYYY HH24:MI'), ''),
       COALESCE(to_char(ends_at, 'MM-DD-YYYY HH24:MI'), ''),
       COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), '') from contracts where id = $1 `,
		id).Scan(
		&contact.ID,
		&contact.CompanyId,
		&contact.Status,
		&contact.StartsAt,
		&contact.EndsAt,
		&contact.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("query contract by id %s: %w", id, err)
	}

	return contact, nil
}

func (r *RepoPg) Update(ctx context.Context, name string, newName string) error {
	_, err := r.pool.Exec(
		ctx,
		"UPDATE contracts SET name = $2 WHERE name = $1",
		name, newName,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) UpdateStatus(ctx context.Context, tx pgx.Tx, id string, newStatus string) error {
	_, err := tx.Exec(ctx, "UPDATE contracts SET status = $2 WHERE id = $1", id, newStatus)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) Delete(ctx context.Context, name string) error {
	_, err := r.pool.Exec(
		ctx,
		"DELETE contracts WHERE name = $1",
		name,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) GetCount(ctx context.Context) (string, error) {
	var count string
	err := r.pool.QueryRow(
		ctx,
		`select count(*) from contracts`,
	).Scan(&count)

	return count, err
}

func (r *RepoPg) DoPing(ctx context.Context) error {
	err := r.pool.Ping(ctx)
	return err
}

func (r *RepoPg) GetForUpdateByID(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*dto.Contract, error) {
	contract := &dto.Contract{}

	err := tx.QueryRow(ctx, "SELECT "+
		"id, "+
		"company_id, "+
		"status, "+
		"COALESCE(to_char(starts_at, 'MM-DD-YYYY HH24:MI'), '') AS starts_at, "+
		"COALESCE(to_char(ends_at, 'MM-DD-YYYY HH24:MI'), '') AS ends_at, "+
		"COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), '') AS created_at "+
		"FROM contracts WHERE id = $1 FOR UPDATE", id).Scan(
		&contract.ID,
		&contract.CompanyId,
		&contract.Status,
		&contract.StartsAt,
		&contract.EndsAt,
		&contract.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContractNotFound
		}
		return nil, fmt.Errorf("query contract by id %s: %w", id, err)
	}

	return contract, nil
}

func (r *RepoPg) CreateContractService(ctx context.Context, it dto.ContractService) (*dto.ContractService, error) {
	_, err := r.pool.Exec(
		ctx,
		`insert into contract_services(id, contract_id, service_code, is_available) values($1, $2, $3, $4)`,
		it.ID, it.ContractId, it.ServiceCode, it.IsAvailable,
	)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Exec: %w", err)
	}

	return &it, nil
}

func (r *RepoPg) ContractServiceList(ctx context.Context, contractId string) ([]dto.ContractService, error) {
	rows, err := r.pool.Query(ctx, `select id, contract_id, service_code, 
       is_available, COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), ''), from contract_services  WHERE contract_id = $1`, contractId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contractServices []dto.ContractService
	for rows.Next() {
		var item dto.ContractService
		if err := rows.Scan(&item.ID, &item.ContractId, &item.ServiceCode, &item.IsAvailable, &item.CreatedAt); err != nil {
			return nil, err
		}
		contractServices = append(contractServices, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contractServices, nil
}

func (r *RepoPg) ContractActiveServiceList(ctx context.Context, contractId string) ([]dto.ContractService, error) {
	rows, err := r.pool.Query(ctx, `select id, contract_id, service_code, 
       is_available, COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), ''), from contract_services  WHERE contract_id = $1 AND is_available = 'TRUE'`, contractId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contractServices []dto.ContractService
	for rows.Next() {
		var item dto.ContractService
		if err := rows.Scan(&item.ID, &item.ContractId, &item.ServiceCode, &item.IsAvailable, &item.CreatedAt); err != nil {
			return nil, err
		}
		contractServices = append(contractServices, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contractServices, nil
}

func (r *RepoPg) GetAvailabilityCompaniesService(ctx context.Context, tx pgx.Tx, serviceCompany dto.ServiceCompanyRequest) (*dto.AvailbaleServicesCompanyResponse, error) {
	var it dto.AvailbaleServicesCompanyResponse

	err := tx.QueryRow(
		ctx,
		`WITH company_contract AS (
    SELECT c.id, c.status, c.starts_at, c.ends_at
    FROM contracts c
    WHERE c.company_id = $1
      AND c.status IN ('active', 'suspended', 'terminated')
    ORDER BY
        CASE c.status
            WHEN 'active' THEN 1
            WHEN 'suspended' THEN 2
            WHEN 'terminated' THEN 3
        END,
        c.created_at DESC
    LIMIT 1
)
SELECT
    CASE
        WHEN cc.id IS NULL THEN false
        WHEN cc.status = 'suspended' THEN false
        WHEN cc.status = 'terminated' THEN false
        WHEN now() < cc.starts_at THEN false
        WHEN cc.ends_at IS NOT NULL AND now() > cc.ends_at THEN false
        WHEN cs.id IS NULL THEN false
        WHEN cs.is_available = false THEN false
        ELSE true
    END AS allowed,
    CASE
        WHEN cc.id IS NULL THEN '`+dto.ReasonNoActiveContract+`'
        WHEN cc.status = 'suspended' THEN '`+dto.ReasonContractSuspended+`'
        WHEN cc.status = 'terminated' THEN '`+dto.ReasonContractTerminated+`'
        WHEN now() < cc.starts_at THEN '`+dto.ReasonContractNotStarted+`'
        WHEN cc.ends_at IS NOT NULL AND now() > cc.ends_at THEN '`+dto.ReasonContractExpired+`'
        WHEN cs.id IS NULL THEN '`+dto.ReasonServiceNotInContract+`'
        WHEN cs.is_available = false THEN '`+dto.ReasonServiceDisabled+`'
        ELSE ''
    END AS reason
FROM (SELECT 1) AS dummy
LEFT JOIN company_contract cc ON true
LEFT JOIN contract_services cs
    ON cs.contract_id = cc.id
   AND cs.service_code = $2 `,
		serviceCompany.CompanyId, serviceCompany.ServiceCode,
	).Scan(&it.Allowed, &it.Reason)

	return &it, err
}
