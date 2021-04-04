package sentry

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

type Sentry struct {
	client *sentry.Client
}

type Meta struct {
	UserID  string
	Request *http.Request
	IsPanic bool
}

const panicTag = "isPanic"

func New(dsn string) (*Sentry, error) {
	if dsn == "" {
		return nil, errors.New("sentry: empty dsn provided")
	}

	cli, err := sentry.NewClient(sentry.ClientOptions{
		Dsn: dsn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sentry: unable to initialise")
	}

	return &Sentry{
		client: cli,
	}, nil
}

func (s *Sentry) CaptureError(ctx context.Context, err error, meta *Meta) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		s.applyMeta(scope, meta)
		s.client.CaptureException(err, nil, nil)
	})
}

func (s *Sentry) Flush(timeout time.Duration) {
	s.client.Flush(timeout)
}

func (s *Sentry) applyMeta(scope *sentry.Scope, meta *Meta) {
	isPanic := false

	if meta != nil {
		if meta.Request != nil {
			scope.SetRequest(meta.Request)
		}
		if meta.UserID != "" {
			scope.SetUser(sentry.User{ID: meta.UserID})
		}

		isPanic = meta.IsPanic
	}

	scope.SetTag(panicTag, strconv.FormatBool(isPanic))
}
