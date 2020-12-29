package slack

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	"github.com/alexsniffin/go-blog2-example/internal/example/models"
)

type Client struct {
	cfg    models.SlackClientConfig
	logger zerolog.Logger
}

func NewClient(cfg models.SlackClientConfig, logger zerolog.Logger) (*Client, error) {
	if cfg.Webhook == "" {
		return nil, errors.New("missing webhook")
	}
	return &Client{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (c *Client) PostMessage(text string) error {
	attachment := slack.Attachment{
		AuthorName: "Example Service",
		Text:       text,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(c.cfg.Webhook, &msg)
	if err != nil {
		return errors.Wrap(err, "failed to post slack message")
	}

	c.logger.Debug().Msgf("posted slack message: %s", text)
	return nil
}
