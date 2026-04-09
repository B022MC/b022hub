//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/B022MC/b022hub/internal/config"
	"github.com/stretchr/testify/require"
)

type rateLimitAccountRepoStub struct {
	mockAccountRepoForGemini
	setErrorCalls int
	tempCalls     int
	lastErrorMsg  string
	bindCalls     int
	bindErr       error
	boundAccount  int64
	boundGroupIDs []int64
}

func (r *rateLimitAccountRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	r.setErrorCalls++
	r.lastErrorMsg = errorMsg
	return nil
}

func (r *rateLimitAccountRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	r.tempCalls++
	return nil
}

func (r *rateLimitAccountRepoStub) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	r.bindCalls++
	r.boundAccount = accountID
	r.boundGroupIDs = append([]int64(nil), groupIDs...)
	return r.bindErr
}

type tokenCacheInvalidatorRecorder struct {
	accounts []*Account
	err      error
}

func (r *tokenCacheInvalidatorRecorder) InvalidateToken(ctx context.Context, account *Account) error {
	r.accounts = append(r.accounts, account)
	return r.err
}

func TestRateLimitService_HandleUpstreamError_OAuth401SetsTempUnschedulable(t *testing.T) {
	t.Run("gemini", func(t *testing.T) {
		repo := &rateLimitAccountRepoStub{}
		invalidator := &tokenCacheInvalidatorRecorder{}
		service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
		service.SetTokenCacheInvalidator(invalidator)
		account := &Account{
			ID:       100,
			Platform: PlatformGemini,
			Type:     AccountTypeOAuth,
			Credentials: map[string]any{
				"temp_unschedulable_enabled": true,
				"temp_unschedulable_rules": []any{
					map[string]any{
						"error_code":       401,
						"keywords":         []any{"unauthorized"},
						"duration_minutes": 30,
						"description":      "custom rule",
					},
				},
			},
		}

		shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, []byte("unauthorized"))

		require.True(t, shouldDisable)
		require.Equal(t, 0, repo.setErrorCalls)
		require.Equal(t, 1, repo.tempCalls)
		require.Len(t, invalidator.accounts, 1)
	})

	t.Run("antigravity_401_uses_SetError", func(t *testing.T) {
		// Antigravity 401 由 applyErrorPolicy 的 temp_unschedulable_rules 控制，
		// HandleUpstreamError 中走 SetError 路径。
		repo := &rateLimitAccountRepoStub{}
		invalidator := &tokenCacheInvalidatorRecorder{}
		service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
		service.SetTokenCacheInvalidator(invalidator)
		account := &Account{
			ID:       100,
			Platform: PlatformAntigravity,
			Type:     AccountTypeOAuth,
		}

		shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, []byte("unauthorized"))

		require.True(t, shouldDisable)
		require.Equal(t, 1, repo.setErrorCalls)
		require.Equal(t, 0, repo.tempCalls)
		require.Empty(t, invalidator.accounts)
	})
}

func TestRateLimitService_HandleUpstreamError_OAuth401InvalidatorError(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	invalidator := &tokenCacheInvalidatorRecorder{err: errors.New("boom")}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	service.SetTokenCacheInvalidator(invalidator)
	account := &Account{
		ID:       101,
		Platform: PlatformGemini,
		Type:     AccountTypeOAuth,
	}

	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, []byte("unauthorized"))

	require.True(t, shouldDisable)
	require.Equal(t, 0, repo.setErrorCalls)
	require.Equal(t, 1, repo.tempCalls)
	require.Len(t, invalidator.accounts, 1)
}

func TestRateLimitService_HandleUpstreamError_NonOAuth401(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	invalidator := &tokenCacheInvalidatorRecorder{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	service.SetTokenCacheInvalidator(invalidator)
	account := &Account{
		ID:       102,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
	}

	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, []byte("unauthorized"))

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Empty(t, invalidator.accounts)
}

func TestRateLimitService_HandleUpstreamError_OpenAIAccountDeactivatedOAuthMovesAccountToUngrouped(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       103,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
	}

	body := []byte(`{"error":{"message":"Your OpenAI account has been deactivated, please check your email for more information.","type":"invalid_request_error","code":"account_deactivated","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, int64(103), repo.boundAccount)
	require.Empty(t, repo.boundGroupIDs)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Contains(t, repo.lastErrorMsg, "OpenAI account deactivated (401)")
}

func TestRateLimitService_HandleUpstreamError_OpenAIAccountDeactivatedAPIKeyMovesAccountToUngrouped(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       104,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
	}

	body := []byte(`{"error":{"message":"Your OpenAI account has been deactivated, please check your email for more information.","type":"invalid_request_error","code":"account_deactivated","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, int64(104), repo.boundAccount)
	require.Empty(t, repo.boundGroupIDs)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Contains(t, repo.lastErrorMsg, "OpenAI account deactivated (401)")
}

func TestRateLimitService_HandleUpstreamError_OpenAITokenRevokedOAuthMovesAccountToUngrouped(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       1041,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
	}

	body := []byte(`{"error":{"message":"Encountered invalidated oauth token for user, failing request","type":null,"code":"token_revoked","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, int64(1041), repo.boundAccount)
	require.Empty(t, repo.boundGroupIDs)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Contains(t, repo.lastErrorMsg, "OpenAI OAuth token revoked (401)")
}

func TestRateLimitService_HandleUpstreamError_OpenAITokenRevokedOAuthUngroupFailureSetsError(t *testing.T) {
	repo := &rateLimitAccountRepoStub{bindErr: errors.New("boom")}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       1042,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
	}

	body := []byte(`{"error":{"message":"Encountered invalidated oauth token for user, failing request","type":null,"code":"token_revoked","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Contains(t, repo.lastErrorMsg, "OpenAI OAuth token revoked (401)")
}

func TestRateLimitService_HandleUpstreamError_OpenAITokenInvalidatedOAuthMovesAccountToUngrouped(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       1043,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
	}

	body := []byte(`{"error":{"message":"Your authentication token has been invalidated. Please try signing in again.","type":"invalid_request_error","code":"token_invalidated","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, int64(1043), repo.boundAccount)
	require.Empty(t, repo.boundGroupIDs)
	require.Equal(t, 1, repo.setErrorCalls)
}

func TestRateLimitService_HandleUpstreamError_SoraTokenInvalidatedMovesAccountToUngrouped(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       1044,
		Platform: PlatformSora,
		Type:     AccountTypeOAuth,
	}

	body := []byte(`{"error":{"message":"Token invalid","code":"token_invalidated"},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.True(t, shouldDisable)
	require.Equal(t, 1, repo.bindCalls)
	require.Equal(t, int64(1044), repo.boundAccount)
	require.Empty(t, repo.boundGroupIDs)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 1, repo.setErrorCalls)
	require.Contains(t, repo.lastErrorMsg, "Sora token invalidated (401)")
}

func TestRateLimitService_HandleUpstreamError_OpenAIAccountDeactivatedPoolModeSkipsLocalState(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       105,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"pool_mode": true,
		},
	}

	body := []byte(`{"error":{"message":"Your OpenAI account has been deactivated, please check your email for more information.","type":"invalid_request_error","code":"account_deactivated","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.False(t, shouldDisable)
	require.Equal(t, 0, repo.bindCalls)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 0, repo.setErrorCalls)
}

func TestRateLimitService_HandleUpstreamError_OpenAIAccountDeactivatedCustomErrorCodesSkip401(t *testing.T) {
	repo := &rateLimitAccountRepoStub{}
	service := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)
	account := &Account{
		ID:       106,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"custom_error_codes_enabled": true,
			"custom_error_codes":         []any{float64(429)},
		},
	}

	body := []byte(`{"error":{"message":"Your OpenAI account has been deactivated, please check your email for more information.","type":"invalid_request_error","code":"account_deactivated","param":null},"status":401}`)
	shouldDisable := service.HandleUpstreamError(context.Background(), account, 401, http.Header{}, body)

	require.False(t, shouldDisable)
	require.Equal(t, 0, repo.bindCalls)
	require.Equal(t, 0, repo.tempCalls)
	require.Equal(t, 0, repo.setErrorCalls)
}
