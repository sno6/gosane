package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Payload struct {
	Text string `json:"text"`
}

type Slack struct {
	webhooks map[string]string
}

func New() *Slack {
	return &Slack{webhooks: make(map[string]string)}
}

func (s *Slack) AddWebhooks(webhooks map[string]string) {
	for key, value := range webhooks {
		s.webhooks[key] = value
	}
}

func (s *Slack) AddWebhook(channel, webhook string) {
	s.webhooks[channel] = webhook
}

func (s *Slack) SendMessage(channel string, message string) error {
	c, ok := s.webhooks[channel]
	if !ok {
		return errors.New("channel not found")
	}

	return s.sendMessageToWebhook(c, message)
}

func (s *Slack) sendMessageToWebhook(webhook string, message string) error {
	payload, err := json.Marshal(Payload{
		Text: message,
	})
	if err != nil {
		return err
	}

	body := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPost, webhook, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
