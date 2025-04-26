package apiclients

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/NarthurN/FIOapi/internal/apiclients/age"
	"github.com/NarthurN/FIOapi/internal/apiclients/gender"
	"github.com/NarthurN/FIOapi/internal/apiclients/nationality"
)

type EnrichmentData struct {
	Age         int
	Gender      string
	Nationality string
}

type Client struct {
	ageClient         *age.Client
	genderClient      *gender.Client
	nationalityClient *nationality.Client
	log               *slog.Logger
}

func New(ageURL, genderURL, nationalityURL string, log *slog.Logger) *Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Client{
		ageClient:         age.New(ageURL, client),
		genderClient:      gender.New(genderURL, client),
		nationalityClient: nationality.New(nationalityURL, client),
		log:               log,
	}
}

func (c *Client) EnrichUserData(ctx context.Context, name string) (*EnrichmentData, error) {
	op := "internal/apiclients/nationality/nationality.go.GetNationality"
	data := &EnrichmentData{}
	var err error

	data.Age, err = c.ageClient.GetAge(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data.Gender, err = c.genderClient.GetGender(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data.Nationality, err = c.nationalityClient.GetNationality(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return data, nil
}
