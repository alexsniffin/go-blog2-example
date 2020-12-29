package server

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-blog2-example/internal/example/clients/prometheus"
	"github.com/alexsniffin/go-blog2-example/internal/example/clients/slack"
	"github.com/alexsniffin/go-blog2-example/internal/example/models"
	"github.com/alexsniffin/go-blog2-example/internal/example/processes/evaluator"
)

type Server struct {
	cfg           models.Config
	logger        zerolog.Logger
	evaluatorPool *evaluator.Pool
	shutdown      sync.Once
}

func NewServer(cfg models.Config, logger zerolog.Logger) (*Server, error) {
	newPrometheusClient, err := prometheus.NewClient(cfg.Prometheus, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create prometheus client")
	}
	newSlackClient, err := slack.NewClient(cfg.Slack, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create slack client")
	}
	newEvaluatorPool, err := evaluator.NewPool(cfg.Evaluator, logger, newPrometheusClient, newSlackClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pool")
	}

	return &Server{
		cfg:           cfg,
		logger:        logger,
		evaluatorPool: newEvaluatorPool,
	}, nil
}

func (s *Server) Start() {
	go s.evaluatorPool.Start()
}

func (s *Server) Shutdown() {
	s.shutdown.Do(s.evaluatorPool.Shutdown)
}
