package user

import (
	"context"
	"time"

	"github.com/sno6/gosane/internal/types"

	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/user"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	client *ent.Client
}

func NewUserStore(client *ent.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.client.User.
		Query().
		Where(
			user.EmailEQ(email),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
}

func (s *Store) FindByUUID(ctx context.Context, uuid uuid.UUID) (*ent.User, error) {
	return s.client.User.
		Query().
		Where(
			user.UUIDEQ(uuid),
			user.DeletedAtIsNil(),
		).
		Only(ctx)
}

func (s *Store) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	providerType := &u.ProviderType
	if u.ProviderType == "" {
		providerType = nil
	}

	var hashedPassword *string
	if u.Password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		hashedPassword = types.String(string(pass))
	}

	return s.client.User.
		Create().
		SetEmail(u.Email).
		SetProviderID(u.ProviderID).
		SetNillableProviderType(providerType).
		SetNillablePassword(hashedPassword).
		SetEmailVerified(u.EmailVerified).
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		Save(ctx)
}

func (s *Store) UpdateByUUID(ctx context.Context, uuid uuid.UUID, u *ent.User) (*ent.User, error) {
	s.client.User.
		Update().
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		Where(
			user.UUIDEQ(uuid),
			user.DeletedAtIsNil(),
		).
		Save(ctx)

	return s.FindByUUID(ctx, uuid)
}

func (s *Store) DeleteByUuid(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.client.User.
		Update().
		SetDeletedAt(time.Now()).
		Where(
			user.UUIDEQ(uuid),
			user.DeletedAtIsNil(),
		).Save(ctx)

	return err
}
