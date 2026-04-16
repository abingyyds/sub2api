package service

import (
	"context"
	"log"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// AdjustPoolBalance 调整分站池余额并写入流水。调用方已保证权限与金额合法性。
// deltaFen > 0 入账，< 0 出账；返回变动后余额。
func (s *SubSiteService) AdjustPoolBalance(ctx context.Context, siteID int64, deltaFen int64, entry SubSiteLedgerEntry) (int64, error) {
	if siteID <= 0 {
		return 0, ErrSubSiteNotFound
	}
	if deltaFen == 0 {
		return 0, infraerrors.BadRequest("SUBSITE_LEDGER_DELTA_ZERO", "ledger delta must be non-zero")
	}
	entry.SubSiteID = siteID
	balance, err := s.repo.AdjustBalance(ctx, siteID, deltaFen, entry)
	if err != nil {
		return 0, err
	}
	if balance < 0 {
		log.Printf("[SubSitePool] sub_site %d balance went negative after %s: %d", siteID, entry.TxType, balance)
	}
	s.invalidateCaches()
	return balance, nil
}

// CreditPoolFromPayment 线上支付成功后将金额入账到分站池（tx_type=topup_online）。
func (s *SubSiteService) CreditPoolFromPayment(ctx context.Context, siteID int64, amountFen int64, orderID int64) error {
	if siteID <= 0 || amountFen <= 0 {
		return infraerrors.BadRequest("SUBSITE_TOPUP_INVALID", "sub-site topup payload is invalid")
	}
	relatedOrder := orderID
	_, err := s.AdjustPoolBalance(ctx, siteID, amountFen, SubSiteLedgerEntry{
		TxType:         SubSiteLedgerTopupOnline,
		RelatedOrderID: &relatedOrder,
		Note:           "线上充值",
	})
	return err
}

// AdminTopupPool 平台管理员给分站池人工加余额（tx_type=topup_admin）。
func (s *SubSiteService) AdminTopupPool(ctx context.Context, siteID int64, amountFen int64, operatorID int64, note string) (*SubSite, error) {
	if amountFen <= 0 {
		return nil, infraerrors.BadRequest("SUBSITE_TOPUP_AMOUNT_INVALID", "topup amount must be positive")
	}
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	entry := SubSiteLedgerEntry{
		TxType: SubSiteLedgerTopupAdmin,
		Note:   strings.TrimSpace(note),
	}
	if operatorID > 0 {
		op := operatorID
		entry.OperatorID = &op
	}
	if _, err := s.AdjustPoolBalance(ctx, siteID, amountFen, entry); err != nil {
		return nil, err
	}
	refreshed, err := s.repo.GetByID(ctx, site.ID)
	if err != nil {
		return nil, err
	}
	return s.populateComputedFields(ctx, refreshed, true)
}

// OfflineTopupUser 分站主给分站内用户线下加余额：user.balance += amount, sub_site.balance -= amount。
// 前置校验：操作人是分站 owner、该 user 已绑定到本分站、allow_offline_topup=true、池余额足够。
func (s *SubSiteService) OfflineTopupUser(ctx context.Context, ownerUserID, siteID, targetUserID int64, amountFen int64, note string) (*SubSite, error) {
	if amountFen <= 0 {
		return nil, infraerrors.BadRequest("SUBSITE_TOPUP_AMOUNT_INVALID", "topup amount must be positive")
	}
	if targetUserID <= 0 {
		return nil, infraerrors.BadRequest("SUBSITE_TOPUP_USER_INVALID", "target user is required")
	}
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site.OwnerUserID != ownerUserID {
		return nil, ErrSubSiteForbidden
	}
	if !site.AllowOfflineTopup {
		return nil, infraerrors.Forbidden("SUBSITE_OFFLINE_TOPUP_DISABLED", "offline topup is disabled for this sub-site")
	}
	if site.BalanceFen < amountFen {
		return nil, infraerrors.BadRequest("SUBSITE_POOL_INSUFFICIENT", "sub-site pool balance is insufficient")
	}
	bound, err := s.repo.GetBoundSubSiteByUserID(ctx, targetUserID)
	if err != nil || bound == nil || bound.ID != site.ID {
		return nil, infraerrors.BadRequest("SUBSITE_USER_NOT_IN_SITE", "target user is not bound to this sub-site")
	}
	if s.userRepo == nil {
		return nil, infraerrors.ServiceUnavailable("USER_REPOSITORY_UNAVAILABLE", "user repository is unavailable")
	}
	// 用户加余额（amount 单位：元）
	if err := s.userRepo.UpdateBalance(ctx, targetUserID, float64(amountFen)/100.0); err != nil {
		return nil, err
	}
	// 从池扣等额
	related := targetUserID
	entry := SubSiteLedgerEntry{
		TxType:        SubSiteLedgerOfflineUserTopup,
		RelatedUserID: &related,
		Note:          strings.TrimSpace(note),
	}
	op := ownerUserID
	entry.OperatorID = &op
	if _, err := s.AdjustPoolBalance(ctx, siteID, -amountFen, entry); err != nil {
		// 线下加用户余额已成功，此处扣池失败留 warn（与 Step 4 扣费路径对齐的容错策略）
		log.Printf("[SubSitePool] offline topup pool debit failed, user=%d site=%d amount=%d: %v", targetUserID, siteID, amountFen, err)
		return nil, err
	}
	refreshed, err := s.repo.GetByID(ctx, site.ID)
	if err != nil {
		return nil, err
	}
	return s.populateComputedFields(ctx, refreshed, true)
}

// DebitPoolForConsumption 分站用户消费时沿 parent 链逐级从分站池扣 1× 基础成本。
// 调用方已确认该用户绑定在 leafSiteID；amountUSD 为用户实际扣费（已经 × 复合倍率）；
// compoundRate 为沿链累乘后的总倍率（用于反推基础成本）。
// 每一个祖先池都扣 1× 基础成本，代表该层向平台购买的服务消耗。
// 允许池透支（负余额）：与 userRepo.DeductBalance 的容忍策略对齐，仅记 warn。
// 调用失败只打日志、不向上游冒泡，避免影响网关主路径。
func (s *SubSiteService) DebitPoolForConsumption(ctx context.Context, leafSiteID, userID, usageLogID int64, amountUSD, compoundRate float64) {
	if leafSiteID <= 0 || compoundRate <= 0 || amountUSD <= 0 {
		return
	}
	basePartFen := int64((amountUSD / compoundRate) * 100.0)
	if basePartFen <= 0 {
		return
	}
	chain, err := s.GetSiteChain(ctx, leafSiteID)
	if err != nil || len(chain) == 0 {
		log.Printf("[SubSitePool] consume chain lookup failed site=%d: %v", leafSiteID, err)
		return
	}
	for _, site := range chain {
		entry := SubSiteLedgerEntry{TxType: SubSiteLedgerConsume}
		if userID > 0 {
			u := userID
			entry.RelatedUserID = &u
		}
		if usageLogID > 0 {
			ul := usageLogID
			entry.RelatedUsageLogID = &ul
		}
		if _, err := s.AdjustPoolBalance(ctx, site.ID, -basePartFen, entry); err != nil {
			log.Printf("[SubSitePool] consume debit failed site=%d user=%d amount=%d: %v", site.ID, userID, basePartFen, err)
		}
	}
}

// GetSiteChain 返回从 leafSiteID 到根分站（含自己）的链，按 leaf→root 顺序。
// 以 MaxSubSiteLevelHardLimit 作为遍历安全上界，防止数据层循环导致死循环。
func (s *SubSiteService) GetSiteChain(ctx context.Context, leafSiteID int64) ([]*SubSite, error) {
	if leafSiteID <= 0 || s == nil || s.repo == nil {
		return nil, nil
	}
	chain := make([]*SubSite, 0, 4)
	seen := make(map[int64]struct{}, 4)
	cursorID := leafSiteID
	for hop := 0; hop <= MaxSubSiteLevelHardLimit; hop++ {
		if _, dup := seen[cursorID]; dup {
			break
		}
		site, err := s.repo.GetByID(ctx, cursorID)
		if err != nil {
			return nil, err
		}
		seen[cursorID] = struct{}{}
		chain = append(chain, site)
		if site.ParentSubSiteID == nil || *site.ParentSubSiteID <= 0 {
			break
		}
		cursorID = *site.ParentSubSiteID
	}
	return chain, nil
}

// CompoundConsumeRateForChain 计算 leaf→root 链上的复合倍率（各层 consume_rate_multiplier 累乘）。
// 任一祖先链中的分站 status != active 则返回 (chain, 0, false)，表示链路已失效、不应按分站计费。
func (s *SubSiteService) CompoundConsumeRateForChain(chain []*SubSite) float64 {
	if len(chain) == 0 {
		return DefaultSubSiteConsumeRate
	}
	rate := 1.0
	for _, site := range chain {
		if site == nil {
			continue
		}
		rate *= normalizeConsumeRateMultiplier(site.ConsumeRateMultiplier)
	}
	if rate <= 0 {
		return DefaultSubSiteConsumeRate
	}
	return rate
}

// ChainActive 判定链路中所有分站是否均处于 active 状态。
// 若链中任一分站被停用，上层应视该用户为"分站已停用"，停止记账。
func (s *SubSiteService) ChainActive(chain []*SubSite) bool {
	for _, site := range chain {
		if site == nil || site.Status != SubSiteStatusActive {
			return false
		}
	}
	return len(chain) > 0
}

// boundSubSiteCacheTTL 控制网关扣费路径查分站绑定的缓存时长。
const boundSubSiteCacheTTL = time.Minute

type boundSubSiteEntry struct {
	site      *SubSite
	expiresAt time.Time
}

// GetBoundSubSiteForUser 返回用户绑定的分站（带 TTL 缓存）。未绑定返回 (nil, nil)。
// 网关扣费路径每次请求都会调用，因此必须廉价。
func (s *SubSiteService) GetBoundSubSiteForUser(ctx context.Context, userID int64) (*SubSite, error) {
	if s == nil || s.repo == nil || userID <= 0 {
		return nil, nil
	}
	now := time.Now()
	s.boundCacheMu.RLock()
	if entry, ok := s.boundCache[userID]; ok && now.Before(entry.expiresAt) {
		s.boundCacheMu.RUnlock()
		return entry.site, nil
	}
	s.boundCacheMu.RUnlock()

	site, err := s.repo.GetBoundSubSiteByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.boundCacheMu.Lock()
	s.boundCache[userID] = boundSubSiteEntry{site: site, expiresAt: now.Add(boundSubSiteCacheTTL)}
	s.boundCacheMu.Unlock()
	return site, nil
}

// invalidateBoundCache 在绑定关系变化时调用（未来扩展：注册/线下加余额时）。
func (s *SubSiteService) invalidateBoundCache(userID int64) {
	if s == nil {
		return
	}
	s.boundCacheMu.Lock()
	delete(s.boundCache, userID)
	s.boundCacheMu.Unlock()
}

// ListPoolLedger 返回分站池流水，支持 txType 过滤。
func (s *SubSiteService) ListPoolLedger(ctx context.Context, siteID int64, params pagination.PaginationParams, txType string) ([]SubSiteLedgerEntry, *pagination.PaginationResult, error) {
	if siteID <= 0 {
		return nil, nil, ErrSubSiteNotFound
	}
	return s.repo.ListLedger(ctx, siteID, params, txType)
}

// ListOwnedPoolLedger 分站主查询自己分站的流水。
func (s *SubSiteService) ListOwnedPoolLedger(ctx context.Context, ownerUserID, siteID int64, params pagination.PaginationParams, txType string) ([]SubSiteLedgerEntry, *pagination.PaginationResult, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, nil, err
	}
	if site.OwnerUserID != ownerUserID {
		return nil, nil, ErrSubSiteForbidden
	}
	return s.repo.ListLedger(ctx, siteID, params, txType)
}
