//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/B022MC/b022hub/internal/config"
	"github.com/stretchr/testify/require"
)

type oauthUserRepoStub struct {
	*userRepoStub
	userByEmail   *User
	getByEmailErr error
}

func (s *oauthUserRepoStub) GetByEmail(ctx context.Context, email string) (*User, error) {
	if s.getByEmailErr != nil {
		return nil, s.getByEmailErr
	}
	if s.userByEmail == nil {
		return nil, ErrUserNotFound
	}
	return s.userByEmail, nil
}

type refreshTokenCacheStub struct{}

func (s *refreshTokenCacheStub) StoreRefreshToken(ctx context.Context, tokenHash string, data *RefreshTokenData, ttl time.Duration) error {
	return nil
}

func (s *refreshTokenCacheStub) GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenData, error) {
	return nil, ErrRefreshTokenNotFound
}

func (s *refreshTokenCacheStub) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return nil
}

func (s *refreshTokenCacheStub) DeleteUserRefreshTokens(ctx context.Context, userID int64) error {
	return nil
}

func (s *refreshTokenCacheStub) DeleteTokenFamily(ctx context.Context, familyID string) error {
	return nil
}

func (s *refreshTokenCacheStub) AddToUserTokenSet(ctx context.Context, userID int64, tokenHash string, ttl time.Duration) error {
	return nil
}

func (s *refreshTokenCacheStub) AddToFamilyTokenSet(ctx context.Context, familyID string, tokenHash string, ttl time.Duration) error {
	return nil
}

func (s *refreshTokenCacheStub) GetUserTokenHashes(ctx context.Context, userID int64) ([]string, error) {
	return nil, nil
}

func (s *refreshTokenCacheStub) GetFamilyTokenHashes(ctx context.Context, familyID string) ([]string, error) {
	return nil, nil
}

func (s *refreshTokenCacheStub) IsTokenInFamily(ctx context.Context, familyID string, tokenHash string) (bool, error) {
	return false, nil
}

func newAuthServiceForOAuthRegistration(repo *oauthUserRepoStub, settings map[string]string) *AuthService {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			ExpireHour: 1,
		},
		Default: config.DefaultConfig{
			UserBalance:     3.5,
			UserConcurrency: 2,
		},
	}

	var settingService *SettingService
	if settings != nil {
		settingService = NewSettingService(&settingRepoStub{values: settings}, cfg)
	}

	return NewAuthService(
		nil,
		repo,
		nil,
		&refreshTokenCacheStub{},
		cfg,
		settingService,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func TestAuthService_LoginOrRegisterOAuthWithTokenPair_FallsBackToLegacyInvitationSwitch(t *testing.T) {
	repo := &oauthUserRepoStub{userRepoStub: &userRepoStub{nextID: 11}}
	service := newAuthServiceForOAuthRegistration(repo, map[string]string{
		SettingKeyRegistrationEnabled:   "true",
		SettingKeyInvitationCodeEnabled: "true",
	})

	_, _, err := service.LoginOrRegisterOAuthWithTokenPair(context.Background(), "user@test.com", "alice", "")
	require.ErrorIs(t, err, ErrOAuthInvitationRequired)
}

func TestAuthService_LoginOrRegisterOAuthWithTokenPair_UsesDedicatedOAuthSignupSettings(t *testing.T) {
	repo := &oauthUserRepoStub{userRepoStub: &userRepoStub{nextID: 12}}
	service := newAuthServiceForOAuthRegistration(repo, map[string]string{
		SettingKeyRegistrationEnabled:        "false",
		SettingKeyInvitationCodeEnabled:      "false",
		SettingKeyOAuthRegistrationEnabled:   "true",
		SettingKeyOAuthInvitationCodeEnabled: "true",
	})

	_, _, err := service.LoginOrRegisterOAuthWithTokenPair(context.Background(), "user@test.com", "alice", "")
	require.ErrorIs(t, err, ErrOAuthInvitationRequired)
}

func TestAuthService_LoginOrRegisterOAuthWithTokenPair_RespectsOAuthSignupDisable(t *testing.T) {
	repo := &oauthUserRepoStub{userRepoStub: &userRepoStub{nextID: 13}}
	service := newAuthServiceForOAuthRegistration(repo, map[string]string{
		SettingKeyRegistrationEnabled:      "true",
		SettingKeyOAuthRegistrationEnabled: "false",
	})

	_, _, err := service.LoginOrRegisterOAuthWithTokenPair(context.Background(), "user@test.com", "alice", "")
	require.ErrorIs(t, err, ErrOAuthRegDisabled)
}

func TestAuthService_LoginOrRegisterOAuthWithTokenPair_DisabledWhenUserCapReached(t *testing.T) {
	repo := &oauthUserRepoStub{userRepoStub: &userRepoStub{nextID: 14}}
	service := newAuthServiceForOAuthRegistration(repo, map[string]string{
		SettingKeyRegistrationEnabled:      "true",
		SettingKeyOAuthRegistrationEnabled: "true",
		SettingKeyRegistrationUserLimit:    "1",
	})
	service.settingService.SetRegistrationUserCountReader(&registrationUserCountReaderStub{count: 1})

	_, _, err := service.LoginOrRegisterOAuthWithTokenPair(context.Background(), "user@test.com", "alice", "")
	require.ErrorIs(t, err, ErrOAuthRegDisabled)
}

func TestAuthService_LoginOrRegisterOAuthWithTokenPair_AllowsExistingUserWhenOAuthSignupDisabled(t *testing.T) {
	repo := &oauthUserRepoStub{
		userRepoStub: &userRepoStub{},
		userByEmail: &User{
			ID:       99,
			Email:    "user@test.com",
			Username: "alice",
			Role:     RoleUser,
			Status:   StatusActive,
		},
	}
	service := newAuthServiceForOAuthRegistration(repo, map[string]string{
		SettingKeyRegistrationEnabled:      "false",
		SettingKeyOAuthRegistrationEnabled: "false",
	})

	pair, user, err := service.LoginOrRegisterOAuthWithTokenPair(context.Background(), "user@test.com", "alice", "")
	require.NoError(t, err)
	require.NotNil(t, pair)
	require.Equal(t, int64(99), user.ID)
}
