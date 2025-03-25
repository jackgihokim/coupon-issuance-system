package coupon

import (
	"errors"
	"sync"
)

type Coupons struct {
	mu    sync.Mutex
	count uint
	list  []*Coupon
}

// NewCoupons initializes a new Coupons instance with the specified count and pre-allocated list capacity.
func NewCoupons(cnt uint) *Coupons {
	return &Coupons{
		count: cnt,
		list:  make([]*Coupon, 0, cnt),
	}
}

// Add inserts a coupon into the list and decrements the available coupons count. Returns an error if no coupons are available.
func (c *Coupons) Add(coupon *Coupon) error {
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
func (c *Coupons) List() []*Coupon {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list
}
