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
	if deltaFen == 0 && entry.TxType != SubSiteLedgerWithdrawPaid {
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
// 仅 pool 模式分站可用；rate 模式分站无独立进货池，用户充值应走主站通道。
// 前置校验：操作人是分站 owner、该 user 已绑定到本分站、allow_offline_topup=true、池余额足够、mode=pool。
// 为修复 parent 池透支：除 leaf 外，沿链上所有 pool 模式祖先同步扣等额（rate 祖先跳过）。
// rate 祖先意味着该级不存在实体池；加用户余额的"基础成本"在用户后续 API 调用时由平台扣，
// 此处不为 rate 祖先预扣。
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
	if site.Mode != SubSiteModePool {
		return nil, infraerrors.Forbidden("SUBSITE_OFFLINE_TOPUP_MODE_INVALID", "offline topup is only available in pool mode")
	}
	if !site.AllowOfflineTopup {
		return nil, infraerrors.Forbidden("SUBSITE_OFFLINE_TOPUP_DISABLED", "offline topup is disabled for this sub-site")
	}
	if site.BalanceFen < amountFen {
		return nil, ErrSubSitePoolInsufficient
	}
	bound, err := s.repo.GetBoundSubSiteByUserID(ctx, targetUserID)
	if err != nil || bound == nil || bound.ID != site.ID {
		return nil, infraerrors.BadRequest("SUBSITE_USER_NOT_IN_SITE", "target user is not bound to this sub-site")
	}
	if s.userRepo == nil {
		return nil, infraerrors.ServiceUnavailable("USER_REPOSITORY_UNAVAILABLE", "user repository is unavailable")
	}

	chain, err := s.GetSiteChain(ctx, siteID)
	if err != nil || len(chain) == 0 {
		return nil, err
	}
	for idx, ancestor := range chain {
		if ancestor == nil || ancestor.Mode != SubSiteModePool {
			continue
		}
		if ancestor.BalanceFen < amountFen {
			if idx == 0 {
				return nil, ErrSubSitePoolInsufficient
			}
			return nil, infraerrors.BadRequest("SUBSITE_PARENT_POOL_INSUFFICIENT", "parent sub-site pool balance is insufficient")
		}
	}

	entries := makePoolDebitEntries(chain, amountFen, targetUserID, ownerUserID, 0, SubSiteLedgerOfflineUserTopup, strings.TrimSpace(note))
	if err := s.repo.ApplyUserBalanceAndPoolLedger(ctx, targetUserID, float64(amountFen)/100.0, entries); err != nil {
		return nil, err
	}
	s.invalidateCaches()
	refreshed, err := s.repo.GetByID(ctx, site.ID)
	if err != nil {
		return nil, err
	}
	return s.populateComputedFields(ctx, refreshed, true)
}

// CreditUserBalanceWithAutoRestock 原子处理 pool 分站用户线上充值：
// 用户余额入账和分站池自动进货扣款必须同事务完成，任一失败都不落半边账。
func (s *SubSiteService) CreditUserBalanceWithAutoRestock(ctx context.Context, userID, siteID, orderID int64, balanceAmount float64, amountFen int64, orderNo string) error {
	if userID <= 0 {
		return ErrUserNotFound
	}
	if siteID <= 0 || amountFen <= 0 {
		return infraerrors.BadRequest("SUBSITE_AUTO_RESTOCK_INVALID", "auto-restock payload is invalid")
	}
	chain, err := s.GetSiteChain(ctx, siteID)
	if err != nil || len(chain) == 0 {
		return err
	}
	if chain[0] == nil || chain[0].Mode != SubSiteModePool {
		return infraerrors.BadRequest("SUBSITE_AUTO_RESTOCK_MODE_INVALID", "auto-restock is only available in pool mode")
	}
	for idx, site := range chain {
		if site == nil || site.Mode != SubSiteModePool {
			continue
		}
		if site.BalanceFen < amountFen {
			if idx == 0 {
				return ErrSubSitePoolInsufficient
			}
			return infraerrors.BadRequest("SUBSITE_PARENT_POOL_INSUFFICIENT", "parent sub-site pool balance is insufficient")
		}
	}
	entries := makePoolDebitEntries(chain, amountFen, userID, 0, orderID, SubSiteLedgerAutoRestock, "用户线上充值自动进货，订单 "+orderNo)
	if err := s.repo.ApplyUserBalanceAndPoolLedger(ctx, userID, balanceAmount, entries); err != nil {
		return err
	}
	s.invalidateCaches()
	return nil
}

func makePoolDebitEntries(chain []*SubSite, amountFen int64, relatedUserID int64, operatorID int64, relatedOrderID int64, txType string, note string) []SubSiteLedgerEntry {
	entries := make([]SubSiteLedgerEntry, 0, len(chain))
	for idx, site := range chain {
		if site == nil || site.Mode != SubSiteModePool {
			continue
		}
		entryNote := note
		if idx > 0 && entryNote != "" {
			entryNote = "[parent sync] " + entryNote
		}
		entry := SubSiteLedgerEntry{
			SubSiteID: site.ID,
			TxType:    txType,
			DeltaFen:  -amountFen,
			Note:      entryNote,
		}
		if relatedUserID > 0 {
			u := relatedUserID
			entry.RelatedUserID = &u
		}
		if operatorID > 0 {
			op := operatorID
			entry.OperatorID = &op
		}
		if relatedOrderID > 0 {
			o := relatedOrderID
			entry.RelatedOrderID = &o
		}
		entries = append(entries, entry)
	}
	return entries
}

// DebitPoolForConsumption 分站用户消费时按链上每个分站的 mode 分别结算：
//
//	pool 模式：扣 1× 基础成本到分站池（代表该级向上游采购的"进货成本"）；
//	rate 模式：不扣池，而是把本级"加价差额"作为分站主利润入账到同一 balance_fen。
//
// 利润公式（rate 级 i 按 leaf→root 索引，0 = leaf）：
//
//	profit_i = base × (compound_i - compound_{i-1}),  compound_{-1} = 1.0
//
// 其中 compound_i 为 leaf 到 i 级累计倍率乘积。
// 这样从用户余额扣掉的 amountUSD = base × compoundRate 会按级别精确分配：
//   - 纯 rate 链：Σ profit_i = base × (compoundRate - 1)，平台留 base
//   - 含 pool 级：pool 级吃 base 扣池，rate 级按差额入账
//
// 允许池/利润账目透支（负余额），仅记 warn，避免影响网关主路径。
// 失败不向上游冒泡。
func (s *SubSiteService) DebitPoolForConsumption(ctx context.Context, leafSiteID, userID, usageLogID int64, amountUSD, compoundRate float64) {
	if leafSiteID <= 0 || compoundRate <= 0 || amountUSD <= 0 {
		return
	}
	baseUSD := amountUSD / compoundRate
	baseFen := int64(baseUSD * 100.0)
	if baseFen <= 0 {
		return
	}
	chain, err := s.GetSiteChain(ctx, leafSiteID)
	if err != nil || len(chain) == 0 {
		log.Printf("[SubSitePool] consume chain lookup failed site=%d: %v", leafSiteID, err)
		return
	}
	prevCompound := 1.0
	for _, site := range chain {
		if site == nil {
			continue
		}
		siteRate := normalizeConsumeRateMultiplier(site.ConsumeRateMultiplier)
		curCompound := prevCompound * siteRate

		entry := SubSiteLedgerEntry{}
		if userID > 0 {
			u := userID
			entry.RelatedUserID = &u
		}
		if usageLogID > 0 {
			ul := usageLogID
			entry.RelatedUsageLogID = &ul
		}

		switch site.Mode {
		case SubSiteModeRate:
			// 按本级加价差额入账利润
			profitFen := int64((curCompound - prevCompound) * baseUSD * 100.0)
			if profitFen > 0 {
				entry.TxType = SubSiteLedgerProfit
				if _, err := s.AdjustPoolBalance(ctx, site.ID, profitFen, entry); err != nil {
					log.Printf("[SubSitePool] profit credit failed site=%d user=%d amount=%d: %v", site.ID, userID, profitFen, err)
				}
			}
		default:
			// pool 模式：按旧逻辑扣 1× base
			entry.TxType = SubSiteLedgerConsume
			if _, err := s.AdjustPoolBalance(ctx, site.ID, -baseFen, entry); err != nil {
				log.Printf("[SubSitePool] consume debit failed site=%d user=%d amount=%d: %v", site.ID, userID, baseFen, err)
			}
		}

		prevCompound = curCompound
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
