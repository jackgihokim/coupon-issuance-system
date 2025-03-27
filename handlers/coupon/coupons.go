package coupon

import (
	"errors"
	"sync"

	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
)

type Coupons struct {
	count uint32
	mu    sync.Mutex
	list  []*couponv1.Coupon
}

// NewCoupons initializes a new Coupons instance with the specified count and pre-allocated list capacity.
func NewCoupons(cnt uint32) *Coupons {
	return &Coupons{
		count: cnt,
		list:  make([]*couponv1.Coupon, 0, cnt),
	}
}

// Add inserts a coupon into the list and decrements the available coupons count. Returns an error if no coupons are available.
func (c *Coupons) Add(coupon *couponv1.Coupon) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.count == 0 {
		return errors.New("no more coupon")
	}
	c.list = append(c.list, coupon)
	c.count--
	return nil
}

// List returns the current list of coupons. This method ensures thread safety by locking during access.
func (c *Coupons) List() []*couponv1.Coupon {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list
}
