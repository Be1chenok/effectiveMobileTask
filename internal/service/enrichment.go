package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	"github.com/Be1chenok/effectiveMobileTask/internal/domain"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
	"go.uber.org/zap"
)

const keyName = "name"

type Enrichment interface {
	GetAgeByName(ctx context.Context, name string) (int, error)
	GetGenderByName(ctx context.Context, name string) (string, error)
	GetNationalityByName(ctx context.Context, name string) (string, error)
}

type enrichment struct {
	conf   *config.Config
	logger appLogger.Logger
}

func NewEnrichment(conf *config.Config, logger appLogger.Logger) Enrichment {
	return &enrichment{
		conf:   conf,
		logger: logger.With(zap.String("component", "service-enrichment")),
	}
}

func (e enrichment) GetAgeByName(ctx context.Context, name string) (int, error) {
	body, err := makeGetRequest(ctx, e.conf.API.AgeURL, keyName, name)
	if err != nil {
		return 0, err
	}

	var resp apiAgeResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("failed to decode JSON data: %w", err)
	}

	e.logger.Debugf("age: %v", resp.Age)

	if resp.Age < 0 {
		return 0, domain.ErrAgeNotFound
	}

	return resp.Age, nil
}

func (e enrichment) GetGenderByName(ctx context.Context, name string) (string, error) {
	body, err := makeGetRequest(ctx, e.conf.API.GenderURL, keyName, name)
	if err != nil {
		return "", nil
	}

	var resp apiGenderResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("failed to decode JSON data: %w", err)
	}

	e.logger.Debugf("gender: %v", resp.Gender)

	if resp.Gender == "" {
		return "", domain.ErrGenderNotFound
	}

	return resp.Gender, nil
}

func (e enrichment) GetNationalityByName(ctx context.Context, name string) (string, error) {
	body, err := makeGetRequest(ctx, e.conf.API.NationalityURL, keyName, name)
	if err != nil {
		return "", err
	}

	var resp apiNationalityResponse

	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("failed to decode JSON data: %w", err)
	}

	e.logger.Debugf("nationality: %v", resp.Nationality)

	if len(resp.Nationality) == 0 {
		return "", domain.ErrNationalityNotFound
	}

	return resp.Nationality[0].CountryId, nil
}

func makeGetRequest(ctx context.Context, url, key, name string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	query := req.URL.Query()
	query.Add(key, name)
	req.URL.RawQuery = query.Encode()

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
