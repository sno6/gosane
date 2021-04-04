package user

import (
	"context"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqljson"
	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/ent/listing"
	"github.com/sno6/gosane/ent/predicate"
	"github.com/sno6/gosane/ent/schema"
	"github.com/sno6/gosane/ent/search"
	"github.com/sno6/gosane/ent/user"
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
	return s.client.User.
		Create().
		SetEmail(u.Email).
		SetProviderID(u.ProviderID).
		SetProviderType(u.ProviderType).
		SetEmailVerified(u.EmailVerified).
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetNotificationSettings(u.NotificationSettings).
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

func (s *Store) UpdateNotificationSettingsByUUID(ctx context.Context, uuid uuid.UUID, ns schema.NotificationSettings) (*ent.User, error) {
	s.client.User.
		Update().
		SetNotificationSettings(ns).
		Where(
			user.UUIDEQ(uuid),
			user.DeletedAtIsNil(),
		).
		Save(ctx)

	return s.FindByUUID(ctx, uuid)
}

func (s *Store) FindUsersWithSearches(ctx context.Context) ([]*ent.User, error) {
	return s.client.User.
		Query().
		WithSearch(
			// Sub query to only find searches that are active.
			func(query *ent.SearchQuery) {
				query.Where(
					search.DeletedAtIsNil(),
					search.StateEQ(schema.Active),
				)
			},
		).
		Where(user.DeletedAtIsNil()).
		All(ctx)
}

func (s *Store) FindUserWithSearch(ctx context.Context, searchUUID uuid.UUID) (*ent.User, error) {
	return s.client.User.
		Query().
		WithSearch(
			func(query *ent.SearchQuery) {
				query.Where(
					search.UUIDEQ(searchUUID),
					search.DeletedAtIsNil(),
					search.StateEQ(schema.Active),
				)
			},
		).
		Where(
			user.DeletedAtIsNil(),
			user.HasSearchWith(
				search.UUIDEQ(searchUUID),
			),
		).
		First(ctx)
}

// Get all users, their searches and their searches' listings that meet the following criteria:
//
// 1. The user is ready to be notified.. meaning their time period matches that of `timestamp`.
// 2. The listings `created_at` > the last time the user was notified.
func (s *Store) FindUsersWithSearchesForNotificationTimestamp(ctx context.Context, timestamp string) ([]*ent.User, error) {
	return s.client.User.
		Query().
		WithSearch(
			func(query *ent.SearchQuery) {
				query.WithListing()
				query.Where(
					search.StateEQ(schema.Active),
					search.DeletedAtIsNil(),
				)
			},
		).
		Where(
			ShouldNotifyUser(timestamp),
			user.DeletedAtIsNil(),
		).
		Order(ent.Desc(listing.FieldCreatedAt)).
		All(ctx)
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

// 1. Does their time match the timestamp? or do they want to get notified per listing.
// 2. Is their email verified?
func ShouldNotifyUser(timestamp string) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s.Where(sqljson.ValueEQ(user.FieldNotificationSettings, timestamp, sqljson.Path("time"))).
			Or().
			Where(sqljson.ValueEQ(user.FieldNotificationSettings, string(schema.PerListing), sqljson.Path("period"))).
			Where(sqljson.ValueEQ(user.FieldNotificationSettings, true, sqljson.Path("emailVerified")))
	})
}
