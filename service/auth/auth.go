package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sno6/gosane/ent"
	"github.com/sno6/gosane/internal/jwt"
	"github.com/sno6/gosane/internal/verification"
	"github.com/sno6/gosane/service/user"
	"github.com/sno6/gosane/store/token"
	"golang.org/x/crypto/bcrypt"
)

var errLogin = errors.New("username or password is incorrect")

type Service struct {
	jwt          *jwt.Auth
	userService  *user.Service
	tokenStore   *token.Store
	verification *verification.Verification
}

func NewAuthService(
	jwt *jwt.Auth,
	tokenStore *token.Store,
	userService *user.Service,
	verification *verification.Verification,
) *Service {
	return &Service{
		jwt:          jwt,
		tokenStore:   tokenStore,
		userService:  userService,
		verification: verification,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*jwt.TokenInfo, error) {
	u, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("username or password is incorrect")
		}

		return nil, err
	}

	if u.ProviderType != nil || u.ProviderID != "" || u.Password == "" {
		return nil, errors.New("unable to login with email/password using a social account")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, errLogin
		}

		return nil, err
	}

	if !u.EmailVerified {
		return nil, errors.New("email is unverified")
	}

	return s.CreateTokens(ctx, u)
}

func (s *Service) UserExists(ctx context.Context, email string) (bool, error) {
	_, err := s.userService.FindByEmail(ctx, email)
	return ent.IsNotFound(err), err
}

func (s *Service) Register(ctx context.Context, u *ent.User) (*ent.User, error) {
	u, err := s.userService.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	if !u.EmailVerified {
		err = s.verification.SendVerificationEmail(u.Email)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *Service) CreateTokens(ctx context.Context, u *ent.User) (*jwt.TokenInfo, error) {
	err := s.tokenStore.CleanseTokensForUser(ctx, u.UUID)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwt.NewTokens(u.UUID.String())
	if err != nil {
		return nil, err
	}

	_, err = s.tokenStore.Create(ctx, &ent.Token{
		RefreshToken:     tokens.Refresh,
		AccessExpiresAt:  tokens.AccessExpiresAt,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	}, u)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*jwt.TokenInfo, error) {
	found, err := s.tokenStore.FindByRefreshToken(ctx, refreshToken, true)
	if err != nil {
		return nil, err
	}

	u := found.Edges.User
	if u == nil {
		return nil, errors.New("no user associated with token")
	}

	err = s.tokenStore.CleanseTokensForUser(ctx, u.UUID)
	if err != nil {
		return nil, err
	}

	newTokens, err := s.CreateTokens(ctx, u)
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *Service) FindUserByToken(ctx context.Context, token string) (*ent.User, error) {
	claims, err := s.jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}

	return s.userService.FindByUUID(ctx, claims.Identifier)
}

func (s *Service) DeleteUserByUuid(ctx context.Context, userUuid uuid.UUID) error {
	_, err := s.userService.FindByUUID(ctx, userUuid.String())
	if err != nil {
		return err
	}

	err = s.tokenStore.DeleteAllTokensForUser(ctx, userUuid)
	if err != nil {
		return err
	}

	err = s.userService.DeleteByUuid(ctx, userUuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	claims, err := s.jwt.ParseToken(token)
	if err != nil {
		return err
	}

	// If the claims are legitimate verify the email for this user.
	email := claims.Identifier
	u, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if u.EmailVerified {
		return errors.New("email already verified")
	}

	_, err = s.userService.UpdateByUUID(ctx, u.UUID, &ent.User{EmailVerified: true})
	return err
}
