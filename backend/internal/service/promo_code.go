package service

import (
	"time"
)

// PromoCode 购买优惠码（折扣/满减）
type PromoCode struct {
	ID             int64
	Code           string
	DiscountAmount float64 // fixed: 减免金额(分), percentage: 百分比值(如10=减10%)
	DiscountType   string  // "fixed" 或 "percentage"
	MinOrderAmount int     // 最低订单金额(分)，0=无门槛
	MaxUses        int
	UsedCount      int
	Status         string
	ExpiresAt      *time.Time
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// 关联
	UsageRecords []PromoCodeUsage
}

// PromoCodeUsage 优惠码使用记录
type PromoCodeUsage struct {
	ID             int64
	PromoCodeID    int64
	UserID         int64
	DiscountAmount float64 // 实际折扣金额(分)
	OrderNo        string  // 关联订单号
	UsedAt         time.Time

	// 关联
	PromoCode *PromoCode
	User      *User
}

// CanUse 检查优惠码是否可用
func (p *PromoCode) CanUse() bool {
	if p.Status != PromoCodeStatusActive {
		return false
	}
	if p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt) {
		return false
	}
	if p.MaxUses > 0 && p.UsedCount >= p.MaxUses {
		return false
	}
	return true
}

// IsExpired 检查是否已过期
func (p *PromoCode) IsExpired() bool {
	return p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt)
}

// CalculateDiscount 根据优惠码类型计算折扣金额（分）
func (p *PromoCode) CalculateDiscount(amountFen int) int {
	if amountFen <= 0 {
		return 0
	}
	var discount int
	switch p.DiscountType {
	case PromoCodeDiscountTypePercentage:
		// DiscountAmount 是百分比值，如 10 表示减 10%
		discount = int(float64(amountFen) * p.DiscountAmount / 100.0)
	default: // "fixed"
		// DiscountAmount 是固定减免金额（分）
		discount = int(p.DiscountAmount)
	}
	// 折扣不能超过原价（至少保留1分）
	if discount >= amountFen {
		discount = amountFen - 1
	}
	if discount < 0 {
		discount = 0
	}
	return discount
}

// CreatePromoCodeInput 创建优惠码输入
type CreatePromoCodeInput struct {
	Code           string
	DiscountAmount float64
	DiscountType   string
	MinOrderAmount int
	MaxUses        int
	ExpiresAt      *time.Time
	Notes          string
}

// UpdatePromoCodeInput 更新优惠码输入
type UpdatePromoCodeInput struct {
	Code           *string
	DiscountAmount *float64
	DiscountType   *string
	MinOrderAmount *int
	MaxUses        *int
	Status         *string
	ExpiresAt      *time.Time
	Notes          *string
}
