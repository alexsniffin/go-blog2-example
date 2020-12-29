package models

type Config struct {
	Logger     Logger
	Evaluator  Evaluator
	Slack      SlackClientConfig
	Prometheus PrometheusClientConfig
}

type Logger struct {
	Level string
}

type Evaluator struct {
	ExprVariableName string
	Rules            []Rule
}

type Rule struct {
	Query       string
	Expression  string
	IntervalSec int
	Template    string
}

type SlackClientConfig struct {
	Webhook string
}

type PrometheusClientConfig struct {
	URL string
}
