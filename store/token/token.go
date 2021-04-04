package token

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/token"
	"github.com/sno6/gosane/ent/user"
)

type Store struct {
	client *ent.Client
}

func NewTokenStore(client *ent.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Create(ctx context.Context, t *ent.Token, u *ent.User) (*ent.Token, error) {
	return s.client.Token.
		Create().
		SetRefreshToken(t.RefreshToken).
		SetRefreshExpiresAt(t.RefreshExpiresAt).
		SetAccessExpiresAt(t.AccessExpiresAt).
		SetUser(u).
		Save(ctx)
}

func (s *Store) FindByRefreshToken(ctx context.Context, refreshToken string, attachUser bool) (*ent.Token, error) {
	q := s.client.Token.
		Query().
		Where(
			token.RefreshTokenEQ(refreshToken),
			token.RefreshExpiresAtGT(time.Now()),
			token.DeletedAtIsNil(),
		)

	if attachUser {
		q.WithUser()
	}

	return q.Only(ctx)
}

func (s *Store) CleanseTokensForUser(ctx context.Context, userUuid uuid.UUID) error {
	_, err := s.client.Token.
		Delete().
		Where(
			token.RefreshExpiresAtLT(time.Now()),
			token.DeletedAtIsNil(),
			token.HasUserWith(
				user.UUIDEQ(userUuid),
			),
		).
		Exec(ctx)

	return err
}

func (s *Store) DeleteAllTokensForUser(ctx context.Context, userUuid uuid.UUID) error {
	_, err := s.client.Token.
		Delete().
		Where(
			token.HasUserWith(
				user.UUIDEQ(userUuid),
			),
		).
		Exec(ctx)

	return err
}
