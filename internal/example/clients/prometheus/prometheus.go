package prometheus

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-blog2-example/internal/example/models"
)

type Client struct {
	v1.API
}

func NewClient(cfg models.PrometheusClientConfig, logger zerolog.Logger) (*Client, error) {
	if cfg.URL == "" {
		return nil, errors.New("missing url")
	}
	client, err := api.NewClient(api.Config{
		Address: cfg.URL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to init prometheus client")
	}

	logger.Info().Msgf("prometheus client init success")
	v1api := v1.NewAPI(client)

	return &Client{v1api}, err
}
