package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type accountExpiryRepoStub struct {
	autoPauseCalls   int
	autoDeleteCalls  int
	autoPauseUpdated int64
	autoDeleteCount  int64
	autoPauseNow     time.Time
	autoDeleteNow    time.Time
	autoPauseErr     error
	autoDeleteErr    error
}

func (s *accountExpiryRepoStub) AutoPauseExpiredAccounts(_ context.Context, now time.Time) (int64, error) {
	s.autoPauseCalls++
	s.autoPauseNow = now
	return s.autoPauseUpdated, s.autoPauseErr
}

func (s *accountExpiryRepoStub) AutoDeleteRateLimitedAccounts(_ context.Context, now time.Time) (int64, error) {
	s.autoDeleteCalls++
	s.autoDeleteNow = now
	return s.autoDeleteCount, s.autoDeleteErr
}

func TestAccountExpiryServiceRunOnce_AutoDeleteDisabled(t *testing.T) {
	repo := &accountExpiryRepoStub{autoPauseUpdated: 2, autoDeleteCount: 3}
	svc := NewAccountExpiryService(repo, time.Second, false)

	svc.runOnce()

	require.Equal(t, 1, repo.autoPauseCalls)
	require.Equal(t, 0, repo.autoDeleteCalls)
	require.False(t, repo.autoPauseNow.IsZero())
}

func TestAccountExpiryServiceRunOnce_AutoDeleteEnabled(t *testing.T) {
	repo := &accountExpiryRepoStub{autoPauseUpdated: 1, autoDeleteCount: 4}
	svc := NewAccountExpiryService(repo, time.Second, true)

	svc.runOnce()

	require.Equal(t, 1, repo.autoPauseCalls)
	require.Equal(t, 1, repo.autoDeleteCalls)
	require.False(t, repo.autoDeleteNow.IsZero())
	require.WithinDuration(t, repo.autoPauseNow, repo.autoDeleteNow, time.Second)
}

func TestAccountExpiryServiceRunOnce_SkipsDeleteWhenPauseFails(t *testing.T) {
	repo := &accountExpiryRepoStub{autoPauseErr: errors.New("pause failed")}
	svc := NewAccountExpiryService(repo, time.Second, true)

	svc.runOnce()

	require.Equal(t, 1, repo.autoPauseCalls)
	require.Equal(t, 0, repo.autoDeleteCalls)
}
