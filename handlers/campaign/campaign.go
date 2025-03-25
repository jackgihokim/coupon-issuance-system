package campaign

import (
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
	"time"
)

type Campaign struct {
	ID          uint
	CouponLimit uint
	Name        string
	Description string
	CreatedAt   time.Time
	StartAt     time.Time
	EndAt       time.Time
	Coupons     *coupon.Coupons
}

// NewCampaign initializes and returns a new Campaign instance with the specified attributes and generated coupons.
func NewCampaign(limit uint, name, desc string, start, end time.Time) *Campaign {
	return &Campaign{
		CouponLimit: limit,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		StartAt:     start,
		EndAt:       end,
		Coupons:     coupon.NewCoupons(limit),
	}
}
