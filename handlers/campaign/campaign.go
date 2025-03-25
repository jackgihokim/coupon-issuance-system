package campaign

import (
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
	"time"
)

type Campaign struct {
	ID          uint
	Limit       uint
	Name        string
	Description string
	CreatedAt   time.Time
	IssueAt     time.Time
	ExpireAt    time.Time
	Coupons     []*coupon.Coupon
}
