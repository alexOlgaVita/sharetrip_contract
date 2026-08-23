package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"job4j.ru/sharetrip-contract/internal/http/dto"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServer_checkServicesAvailability(t *testing.T) {
	t.Run("Проверка доступности услуги - договор в статусе 'active', startAt<= now() <endAt - allowed, disabled, not in Contract", func(t *testing.T) {

		// 1. создаем контракт
		// Это аналог NOW() в Postgres для TIMESTAMPTZ: используем .UTC(), чтобы избежать проблем с локальным временем сервера
		now := time.Now().UTC()
		// Добавляем 1 секунду,чтобы заявка согла создать с началом даты старта "в будущем"
		today := now.Add(time.Second)
		// Добавляем 1 день (0 лет, 0 месяцев, 1 день)
		tomorrow := now.AddDate(0, 0, 1)
		layout := "2006-01-02 15:04:05"
		payload := dto.CreateContractRequest{
			CompanyId: uuid.NewString(),
			StartsAt:  today.Format(layout),
			EndsAt:    tomorrow.Format(layout),
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got dto.Contract
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftContract(t, got, payload)
		contractId := got.ID
		companyId := got.CompanyId

		// 2. переводим контракт в активный статус
		payloadStatus := dto.ChangeContractsStatusModelRequest{
			Status: dto.ContractStatusActive,
		}

		body, err = json.Marshal(payloadStatus)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/contracts/"+contractId+"/status-changes",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		got.Status = dto.ContractStatusActive
		var gotStatus dto.Contract
		err = json.Unmarshal(respBody, &gotStatus)
		require.NoError(t, err)
		requireEqualPublishedContract(t, got, gotStatus)

		// 3. добавляем сервисы к контракту компании
		payloadService1 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.create",
			Enabled:     true,
		}
		payloadService2 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.share",
			Enabled:     false,
		}
		var payloadServicesArr []dto.ContractServiceItemRequest
		payloadServicesArr = append(payloadServicesArr, payloadService1, payloadService2)
		payloadServices := dto.ServicesContractRequest{
			Services: payloadServicesArr,
		}

		body, err = json.Marshal(payloadServices)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/contracts/"+contractId+"/services",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServices []dto.ContractService
		err = json.Unmarshal(respBody, &gotServices)
		require.NoError(t, err)

		payloadAsServices := []dto.ContractService{
			{
				ContractId:  contractId,
				ServiceCode: payloadService1.ServiceCode,
				IsAvailable: payloadService1.Enabled,
			},
			{
				ContractId:  contractId,
				ServiceCode: payloadService2.ServiceCode,
				IsAvailable: payloadService2.Enabled,
			},
		}
		requireEqualServicesContract(t, payloadAsServices, gotServices)

		// 4.1 проверка доступности enabled-service
		// Можно задать задержку, например, на 1 секунду, если тест будет отрабатывать очень быстро
		//и таким обр., время старта останется в будущем
		time.Sleep(1 * time.Second)

		payloadServiceCompany := dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.create",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServiceCompany dto.AvailbaleServicesCompanyResponse
		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanySucces(t, gotServiceCompany)

		// 4.2 проверка доступности disabled-service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.share",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonServiceDisabled)

		// 4.3 проверка доступности отсутсвующего в услугах компании service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.not",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonServiceNotInContract)

	})

	t.Run("Проверка доступности услуги - договор в статусе 'draft', startAt<= now() <endAt - allowed, disabled, not in Contract", func(t *testing.T) {

		// 1. создаем контракт
		// Это аналог NOW() в Postgres для TIMESTAMPTZ: используем .UTC(), чтобы избежать проблем с локальным временем сервера
		now := time.Now().UTC()
		// Добавляем 1 секунду,чтобы заявка согла создать с началом даты старта "в будущем"
		today := now.Add(time.Second)
		// Добавляем 1 день (0 лет, 0 месяцев, 1 день)
		tomorrow := now.AddDate(0, 0, 1)
		layout := "2006-01-02 15:04:05"
		payload := dto.CreateContractRequest{
			CompanyId: uuid.NewString(),
			StartsAt:  today.Format(layout),
			EndsAt:    tomorrow.Format(layout),
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got dto.Contract
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftContract(t, got, payload)
		contractId := got.ID
		companyId := got.CompanyId

		// 2. оставляем контакт в статусе "draft"

		// 3. добавляем сервисы к контракту компании
		payloadService1 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.create",
			Enabled:     true,
		}
		payloadService2 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.share",
			Enabled:     false,
		}
		var payloadServicesArr []dto.ContractServiceItemRequest
		payloadServicesArr = append(payloadServicesArr, payloadService1, payloadService2)
		payloadServices := dto.ServicesContractRequest{
			Services: payloadServicesArr,
		}

		body, err = json.Marshal(payloadServices)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/contracts/"+contractId+"/services",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServices []dto.ContractService
		err = json.Unmarshal(respBody, &gotServices)
		require.NoError(t, err)

		payloadAsServices := []dto.ContractService{
			{
				ContractId:  contractId,
				ServiceCode: payloadService1.ServiceCode,
				IsAvailable: payloadService1.Enabled,
			},
			{
				ContractId:  contractId,
				ServiceCode: payloadService2.ServiceCode,
				IsAvailable: payloadService2.Enabled,
			},
		}
		requireEqualServicesContract(t, payloadAsServices, gotServices)

		// 4.1 проверка доступности enabled-service
		// Можно задать задержку, например, на 1 секунду, если тест будет отрабатывать очень быстро
		//и таким обр., время старта останется в будущем
		//time.Sleep(1 * time.Second)

		payloadServiceCompany := dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.create",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServiceCompany dto.AvailbaleServicesCompanyResponse
		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonNoActiveContract)

		// 4.2 проверка доступности disabled-service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.share",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonNoActiveContract)

		// 4.3 проверка доступности отсутсвующего в услугах компании service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.not",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonNoActiveContract)

	})

	t.Run("Проверка доступности услуги - договор в статусе 'active', startAt > now()  - allowed, disabled, not in Contract", func(t *testing.T) {

		// 1. создаем контракт
		// Это аналог NOW() в Postgres для TIMESTAMPTZ: используем .UTC(), чтобы избежать проблем с локальным временем сервера
		now := time.Now().UTC()
		// Добавляем 1 секунду,чтобы заявка согла создать с началом даты старта "в будущем"
		today := now.Add(time.Hour)
		// Добавляем 1 день (0 лет, 0 месяцев, 1 день)
		tomorrow := now.AddDate(0, 0, 1)
		layout := "2006-01-02 15:04:05"
		payload := dto.CreateContractRequest{
			CompanyId: uuid.NewString(),
			StartsAt:  today.Format(layout),
			EndsAt:    tomorrow.Format(layout),
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got dto.Contract
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		requireEqualCreatedDraftContract(t, got, payload)
		contractId := got.ID
		companyId := got.CompanyId

		// 2. переводим контракт в активный статус
		payloadStatus := dto.ChangeContractsStatusModelRequest{
			Status: dto.ContractStatusActive,
		}

		body, err = json.Marshal(payloadStatus)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/contracts/"+contractId+"/status-changes",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		got.Status = dto.ContractStatusActive
		var gotStatus dto.Contract
		err = json.Unmarshal(respBody, &gotStatus)
		require.NoError(t, err)
		requireEqualPublishedContract(t, got, gotStatus)

		// 3. добавляем сервисы к контракту компании
		payloadService1 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.create",
			Enabled:     true,
		}
		payloadService2 := dto.ContractServiceItemRequest{
			ServiceCode: "trip.share",
			Enabled:     false,
		}
		var payloadServicesArr []dto.ContractServiceItemRequest
		payloadServicesArr = append(payloadServicesArr, payloadService1, payloadService2)
		payloadServices := dto.ServicesContractRequest{
			Services: payloadServicesArr,
		}

		body, err = json.Marshal(payloadServices)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPut,
			"/contracts/"+contractId+"/services",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServices []dto.ContractService
		err = json.Unmarshal(respBody, &gotServices)
		require.NoError(t, err)

		payloadAsServices := []dto.ContractService{
			{
				ContractId:  contractId,
				ServiceCode: payloadService1.ServiceCode,
				IsAvailable: payloadService1.Enabled,
			},
			{
				ContractId:  contractId,
				ServiceCode: payloadService2.ServiceCode,
				IsAvailable: payloadService2.Enabled,
			},
		}
		requireEqualServicesContract(t, payloadAsServices, gotServices)

		// 4.1 проверка доступности enabled-service
		// Можно задать задержку, например, на 1 секунду, если тест будет отрабатывать очень быстро
		//и таким обр., время старта останется в будущем
		//time.Sleep(1 * time.Second)

		payloadServiceCompany := dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.create",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var gotServiceCompany dto.AvailbaleServicesCompanyResponse
		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonContractNotStarted)

		// 4.2 проверка доступности disabled-service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.share",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonContractNotStarted)

		// 4.3 проверка доступности отсутсвующего в услугах компании service
		payloadServiceCompany = dto.ServiceCompanyRequest{
			CompanyId:   companyId,
			ServiceCode: "trip.not",
		}

		body, err = json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err = http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = json.Unmarshal(respBody, &gotServiceCompany)
		require.NoError(t, err)
		requireEqualAvailbaleServicesCompanyNotSucces(t, gotServiceCompany, dto.ReasonContractNotStarted)

	})

	t.Run("Проверка доступности услуги - невалидные данные - не задан serviceCode", func(t *testing.T) {

		// 1. не создаем контракт
		// 2. не переводим контракт в активный статус
		// 3. не добавляем сервисы к контракту компании
		// 4 проверка доступности: serviceCode не задан в запросе

		payloadServiceCompany := dto.ServiceCompanyRequest{
			CompanyId: uuid.NewString(),
		}

		body, err := json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"serviceCode is required")
	})

	t.Run("Проверка доступности услуги - невалидные данные - не задан companyId", func(t *testing.T) {

		// 1. не создаем контракт
		// 2. не переводим контракт в активный статус
		// 3. не добавляем сервисы к контракту компании
		// 4 проверка доступности: companyId не задан в запросе
		payloadServiceCompany := dto.ServiceCompanyRequest{
			ServiceCode: "trip.create",
		}

		body, err := json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"companyId is required")
	})

	t.Run("Проверка доступности услуги - невалидные данные - неверный формат companyId", func(t *testing.T) {

		// 1. не создаем контракт
		// 2. не переводим контракт в активный статус
		// 3. не добавляем сервисы к контракту компании
		// 4 проверка доступности: неверный формат companyId

		payloadServiceCompany := dto.ServiceCompanyRequest{
			CompanyId:   "1" + uuid.NewString(),
			ServiceCode: "trip.create",
		}

		body, err := json.Marshal(payloadServiceCompany)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/contracts/service-availability-checks",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			string(respBody),
			"companyId must be a valid UUID")
	})

}

func requireEqualCreatedDraftContract(t *testing.T, got dto.Contract, payload dto.CreateContractRequest) {
	t.Helper()
	require.NotEmpty(t, got.ID)
	require.Equal(t, dto.Contract{
		ID:        got.ID,
		CompanyId: payload.CompanyId,
		Status:    dto.ContractStatusDraft,
		StartsAt:  payload.StartsAt,
		EndsAt:    payload.EndsAt,
		CreatedAt: got.CreatedAt,
	}, got)
}

func requireEqualPublishedContract(t *testing.T, created dto.Contract, got dto.Contract) {
	t.Helper()
	require.NotEmpty(t, got.ID)
	require.NotEmpty(t, got.CreatedAt)

	dateTimeFormatStartsAt, err := dateToUserFormat(created.StartsAt)
	require.NoError(t, err)
	dateTimeFormatEndsAt, err := dateToUserFormat(created.EndsAt)
	require.NoError(t, err)

	require.Equal(t, dto.Contract{
		ID:        created.ID,
		CompanyId: created.CompanyId,
		StartsAt:  dateTimeFormatStartsAt,
		EndsAt:    dateTimeFormatEndsAt,
		Status:    dto.ContractStatusActive,
		CreatedAt: got.CreatedAt,
	}, got)
}

func requireEqualServicesContract(t *testing.T, payloadServices []dto.ContractService, gotServices []dto.ContractService) {
	t.Helper()
	require.Equal(t, len(payloadServices), len(gotServices))

	for i := range gotServices {
		require.NotEmpty(t, gotServices[i].ID)
		gotServices[i].ID = ""
	}
	require.Equal(t, payloadServices, gotServices)
}

func dateToUserFormat(dateTimeIn string) (string, error) {
	templateDate := "2006-01-02 15:04:05"
	dateTimeOutParse, err := time.Parse(templateDate, dateTimeIn)
	if err != nil {
		return "", err
	}
	outputLayout := "01-02-2006 15:04"
	dateTimeOut := dateTimeOutParse.Format(outputLayout)
	return dateTimeOut, nil
}

func requireEqualAvailbaleServicesCompanySucces(t *testing.T, got dto.AvailbaleServicesCompanyResponse) {
	t.Helper()

	require.Equal(t, dto.AvailbaleServicesCompanyResponse{
		Allowed: true,
	}, got)
}

func requireEqualAvailbaleServicesCompanyNotSucces(t *testing.T, got dto.AvailbaleServicesCompanyResponse, reason string) {
	t.Helper()

	require.Equal(t, dto.AvailbaleServicesCompanyResponse{
		Allowed: false,
		Reason:  reason,
	}, got)
}
