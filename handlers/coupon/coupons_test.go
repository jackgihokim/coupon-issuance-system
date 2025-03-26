package coupon

import (
	"testing"

	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
)

func TestNewCoupons(t *testing.T) {
	// Test case for the NewCoupons function
	count := uint64(5)
	coupons := NewCoupons(count)

	if coupons.count != count {
		t.Errorf("Expected count to be %d, got %d", count, coupons.count)
	}

	if cap(coupons.list) != int(count) {
		t.Errorf("Expected list capacity to be %d, got %d", count, cap(coupons.list))
	}

	if len(coupons.list) != 0 {
		t.Errorf("Expected list length to be 0, got %d", len(coupons.list))
	}
}

func TestCoupons_Add(t *testing.T) {
	// Test case for successfully adding coupons
	coupons := NewCoupons(2)
	coupon1 := &couponv1.Coupon{} // Assuming Coupon is defined elsewhere
	coupon2 := &couponv1.Coupon{}

	// First add should succeed
	err := coupons.Add(coupon1)
	if err != nil {
		t.Errorf("Expected no error when adding first coupon, got: %v", err)
	}

	if len(coupons.list) != 1 {
		t.Errorf("Expected list length to be 1 after adding, got %d", len(coupons.list))
	}

	// Second add should succeed
	err = coupons.Add(coupon2)
	if err != nil {
		t.Errorf("Expected no error when adding second coupon, got: %v", err)
	}

	// Third add should fail with "no more coupon" error
	err = coupons.Add(&couponv1.Coupon{})
	if err == nil {
		t.Error("Expected error when adding beyond capacity, got nil")
	}

	if err != nil && err.Error() != "no more coupon" {
		t.Errorf("Expected 'no more coupon' error, got: %v", err)
	}
}

func TestCoupons_List(t *testing.T) {
	// Test case for listing coupons
	coupons := NewCoupons(3)
	coupon1 := &couponv1.Coupon{}
	coupon2 := &couponv1.Coupon{}

	// Add coupons to the list
	_ = coupons.Add(coupon1)
	_ = coupons.Add(coupon2)

	// Get the list and verify
	list := coupons.List()

	if len(list) != 2 {
		t.Errorf("Expected list length to be 2, got %d", len(list))
	}

	if list[0] != coupon1 {
		t.Errorf("Expected first coupon to be %v, got %v", coupon1, list[0])
	}

	if list[1] != coupon2 {
		t.Errorf("Expected second coupon to be %v, got %v", coupon2, list[1])
	}
}

// TestConcurrentAccess tests thread safety of the Coupons methods
func TestConcurrentAccess(t *testing.T) {
	coupons := NewCoupons(100)
	done := make(chan bool)

	// Launch multiple goroutines to add coupons concurrently
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				_ = coupons.Add(&couponv1.Coupon{})
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify the count and list length
	list := coupons.List()
	if len(list) != 100 {
		t.Errorf("Expected list length to be 100, got %d", len(list))
	}

	if coupons.count != 0 {
		t.Errorf("Expected count to be 0, got %d", coupons.count)
	}
}
