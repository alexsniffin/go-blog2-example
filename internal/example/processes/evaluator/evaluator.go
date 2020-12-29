package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-blog2-example/internal/example/clients/prometheus"
	"github.com/alexsniffin/go-blog2-example/internal/example/clients/slack"
	"github.com/alexsniffin/go-blog2-example/internal/example/models"
)

/*
 * Note: For the sake of this example, I clumped everything into here, but ideally the logic should be separated! :)
 */

type Pool struct {
	workers []*worker
	done    chan bool
}

type worker struct {
	logger           zerolog.Logger
	rule             models.Rule
	prometheusClient *prometheus.Client
	slackClient      *slack.Client
	template         *template.Template
	exprVariableName string
	env              map[string]interface{}
	program          *vm.Program
	vm               vm.VM
	ticker           *time.Ticker
	done             chan bool
}

func NewPool(cfg models.Evaluator, logger zerolog.Logger, pClient *prometheus.Client, sClient *slack.Client) (*Pool, error) {
	if len(cfg.Rules) == 0 {
		return nil, errors.New("missing rules")
	}

	done := make(chan bool)

	var workers []*worker
	for _, rule := range cfg.Rules {
		t, err := template.New("message").Parse(rule.Template)
		if err != nil {
			return nil, errors.Wrap(err, "failed to template message")
		}
		env := map[string]interface{}{
			cfg.ExprVariableName: 0,
			"sprintf":            fmt.Sprintf,
		}
		program, err := expr.Compile(rule.Expression, expr.Env(env))
		if err != nil {
			return nil, errors.New("failed to compile expression")
		}
		workers = append(workers, &worker{
			logger:           logger,
			rule:             rule,
			prometheusClient: pClient,
			slackClient:      sClient,
			template:         t,
			exprVariableName: cfg.ExprVariableName,
			env:              env,
			program:          program,
			vm:               vm.VM{},
			done:             done,
		})
	}

	return &Pool{
		workers: workers,
		done:    done,
	}, nil
}

func (p *Pool) Start() {
	for _, wkr := range p.workers {
		go func(w *worker) {
			w.ticker = time.NewTicker(time.Duration(w.rule.IntervalSec) * time.Second)
			w.run()
		}(wkr)
	}
}

func (p *Pool) Shutdown() {
	close(p.done)
}

func (w *worker) run() {
	for {
		select {
		case <-w.done:
			w.ticker.Stop()
			return
		case <-w.ticker.C:
			err := w.evaluate()
			if err != nil {
				w.logger.Error().Err(err).Msg("failure evaluating rule")
			}
		}
	}
}

func (w *worker) evaluate() error {
	res, warnings, err := w.prometheusClient.Query(context.Background(), w.rule.Query, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to query prometheus")
	}
	if len(warnings) > 0 {
		w.logger.Warn().Msgf("prometheus warnings: %v", warnings)
	}

	// match the response to vector and print the response values
	switch r := res.(type) {
	case model.Vector:
		// we can assume for this example the len will be 1, but for a real world problem we'd want to add functionality
		// to sum the values or require the expression to use a map:
		// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md#builtin-functions
		if r.Len() != 1 {
			return errors.New("unexpected result length")
		}
		v, err := strconv.ParseFloat(r[0].Value.String(), 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse response to float")
		}

		w.logger.Debug().Msgf("prometheus result query=%s with value=%d", w.rule.Query, v)

		env := w.env
		env[w.exprVariableName] = v
		// run the compiled expression with the values from env
		out, err := w.vm.Run(w.program, env)
		if err != nil {
			return errors.Wrap(err, "failed to run expression evaluation")
		}

		exprRes := struct {
			Output string
		}{
			Output: fmt.Sprint(out),
		}

		var tpl bytes.Buffer
		if err := w.template.Execute(&tpl, exprRes); err != nil {
			return errors.Wrap(err, "failed to execute template")
		}

		err = w.slackClient.PostMessage(tpl.String())
		if err != nil {
			return err
		}
	default:
		return errors.New("prometheus response type not implemented")
	}
	return nil
}
