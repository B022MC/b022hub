//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/B022MC/b022hub/internal/config"
	"github.com/stretchr/testify/require"
)

func TestSettingService_UpdateSettings_StoresOAuthFlags(t *testing.T) {
	repo := &settingUpdateRepoStub{}
	svc := NewSettingService(repo, &config.Config{})

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		OAuthRegistrationEnabled:   true,
		OAuthInvitationCodeEnabled: true,
	})
	require.NoError(t, err)
	require.Equal(t, "true", repo.updates[SettingKeyOAuthRegistrationEnabled])
	require.Equal(t, "true", repo.updates[SettingKeyOAuthInvitationCodeEnabled])
}

func TestSettingService_ParseSettings_OAuthFlagsFallbackToLegacyRegistration(t *testing.T) {
	svc := NewSettingService(nil, &config.Config{})

	got := svc.parseSettings(map[string]string{
		SettingKeyRegistrationEnabled:   "false",
		SettingKeyInvitationCodeEnabled: "true",
	})
	require.False(t, got.OAuthRegistrationEnabled)
	require.True(t, got.OAuthInvitationCodeEnabled)

	got = svc.parseSettings(map[string]string{
		SettingKeyRegistrationEnabled:        "false",
		SettingKeyInvitationCodeEnabled:      "false",
		SettingKeyOAuthRegistrationEnabled:   "true",
		SettingKeyOAuthInvitationCodeEnabled: "true",
	})
	require.True(t, got.OAuthRegistrationEnabled)
	require.True(t, got.OAuthInvitationCodeEnabled)
}
