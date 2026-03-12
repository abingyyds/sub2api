package service

import (
	"context"
	"log"
	"sync"
	"time"
)

// BalanceExpiryService periodically clears expired initial balances.
type BalanceExpiryService struct {
	userRepo     UserRepository
	billingCache *BillingCacheService
	interval     time.Duration
	stopCh       chan struct{}
	stopOnce     sync.Once
	wg           sync.WaitGroup
}

func NewBalanceExpiryService(userRepo UserRepository, billingCache *BillingCacheService, interval time.Duration) *BalanceExpiryService {
	return &BalanceExpiryService{
		userRepo:     userRepo,
		billingCache: billingCache,
		interval:     interval,
		stopCh:       make(chan struct{}),
	}
}

func (s *BalanceExpiryService) Start() {
	if s == nil || s.userRepo == nil || s.interval <= 0 {
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

func (s *BalanceExpiryService) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		close(s.stopCh)
	})
	s.wg.Wait()
}

func (s *BalanceExpiryService) runOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userIDs, err := s.userRepo.ClearExpiredInitialBalances(ctx)
	if err != nil {
		log.Printf("[BalanceExpiry] Clear expired initial balances failed: %v", err)
		return
	}
	if len(userIDs) > 0 {
		log.Printf("[BalanceExpiry] Cleared expired initial balances for %d users: %v", len(userIDs), userIDs)

		// Invalidate Redis balance cache for affected users
		if s.billingCache != nil {
			for _, uid := range userIDs {
				if err := s.billingCache.InvalidateUserBalance(ctx, uid); err != nil {
					log.Printf("[BalanceExpiry] Failed to invalidate balance cache for user %d: %v", uid, err)
				}
			}
		}
	}
}
