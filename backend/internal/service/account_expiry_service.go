package service

import (
	"context"
	"log"
	"sync"
	"time"
)

type rateLimitedAccountAutoDeleteRepository interface {
	AutoDeleteRateLimitedAccounts(ctx context.Context, now time.Time) (int64, error)
}

type accountExpiryRepository interface {
	AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error)
}

// AccountExpiryService periodically pauses expired accounts and can optionally
// delete rate-limited accounts that were explicitly marked for auto-delete.
type AccountExpiryService struct {
	accountRepo           accountExpiryRepository
	rateLimitedDeleteRepo rateLimitedAccountAutoDeleteRepository
	autoDeleteRateLimited bool
	interval              time.Duration
	stopCh                chan struct{}
	stopOnce              sync.Once
	wg                    sync.WaitGroup
}

func NewAccountExpiryService(accountRepo accountExpiryRepository, interval time.Duration, autoDeleteRateLimited bool) *AccountExpiryService {
	svc := &AccountExpiryService{
		accountRepo:           accountRepo,
		autoDeleteRateLimited: autoDeleteRateLimited,
		interval:              interval,
		stopCh:                make(chan struct{}),
	}
	if deleteRepo, ok := accountRepo.(rateLimitedAccountAutoDeleteRepository); ok {
		svc.rateLimitedDeleteRepo = deleteRepo
	}
	return svc
}

func (s *AccountExpiryService) Start() {
	if s == nil || s.accountRepo == nil || s.interval <= 0 {
		return
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		s.runOnce()
		for {
			select {
			case <-ticker.C:
				s.runOnce()
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *AccountExpiryService) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
	s.wg.Wait()
}

func (s *AccountExpiryService) runOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	updated, err := s.accountRepo.AutoPauseExpiredAccounts(ctx, now)
	if err != nil {
		log.Printf("[AccountExpiry] Auto pause expired accounts failed: %v", err)
		return
	}
	if updated > 0 {
		log.Printf("[AccountExpiry] Auto paused %d expired accounts", updated)
	}

	if !s.autoDeleteRateLimited || s.rateLimitedDeleteRepo == nil {
		return
	}

	deleted, err := s.rateLimitedDeleteRepo.AutoDeleteRateLimitedAccounts(ctx, now)
	if err != nil {
		log.Printf("[AccountExpiry] Auto delete rate-limited accounts failed: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("[AccountExpiry] Auto deleted %d rate-limited accounts", deleted)
	}
}
