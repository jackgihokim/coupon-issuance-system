package coupon

import "time"

type Coupon struct {
	ID       uint
	Code     string
	IssuedAt time.Time
}
