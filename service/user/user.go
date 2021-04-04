package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/internal/merge"
	"github.com/sno6/gosane/store/user"
)

type Service struct {
	userStore *user.Store
}

func NewUserService(userStore *user.Store) *Service {
	return &Service{
		userStore: userStore,
	}
}

func (s *Service) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	return s.userStore.Create(ctx, u)
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.userStore.FindByEmail(ctx, email)
}

func (s *Service) FindByUUID(ctx context.Context, uid string) (*ent.User, error) {
	return s.userStore.FindByUUID(ctx, uuid.MustParse(uid))
}

func (s *Service) UpdateByUUID(ctx context.Context, uuid uuid.UUID, u *ent.User) (*ent.User, error) {
	existing, err := s.FindByUUID(ctx, uuid.String())
	if err != nil {
		return nil, err
	}

	err = merge.Merge(existing, u)
	if err != nil {
		return nil, err
	}

	_, err = s.userStore.UpdateByUUID(ctx, uuid, existing)
	if err != nil {
		return nil, err
	}

	return s.FindByUUID(ctx, uuid.String())
}

func (s *Service) DeleteByUuid(ctx context.Context, userUuid uuid.UUID) error {
	return s.userStore.DeleteByUuid(ctx, userUuid)
}
