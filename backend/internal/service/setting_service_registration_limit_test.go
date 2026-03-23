//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/B022MC/b022hub/internal/config"
	"github.com/stretchr/testify/require"
)

type registrationUserCountReaderStub struct {
	count int64
	err   error
}

func (s *registrationUserCountReaderStub) CountUsers(context.Context) (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.count, nil
}

func TestSettingService_GetPublicSettings_DisablesRegistrationWhenUserCapReached(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyRegistrationEnabled:        "true",
			SettingKeyRegistrationUserLimit:      "2",
			SettingKeyOAuthRegistrationEnabled:   "true",
			SettingKeyInvitationCodeEnabled:      "false",
			SettingKeyOAuthInvitationCodeEnabled: "false",
		},
	}
	svc := NewSettingService(repo, &config.Config{})
	svc.SetRegistrationUserCountReader(&registrationUserCountReaderStub{count: 2})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.False(t, settings.RegistrationEnabled)
	require.False(t, settings.OAuthRegistrationEnabled)
}

func TestSettingService_IsRegistrationEnabled_FailsOpenOnCountError(t *testing.T) {
	svc := NewSettingService(&settingRepoStub{values: map[string]string{
		SettingKeyRegistrationEnabled:   "true",
		SettingKeyRegistrationUserLimit: "1",
	}}, &config.Config{})
	svc.SetRegistrationUserCountReader(&registrationUserCountReaderStub{err: errors.New("db down")})

	require.True(t, svc.IsRegistrationEnabled(context.Background()))
}
