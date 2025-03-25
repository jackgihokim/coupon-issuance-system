package coupon

import (
	"strings"
	"testing"
	"time"
)

func TestNewCoupon(t *testing.T) {
	coupon, err := NewCoupon()
	if err != nil {
		t.Fatalf("Error occurred while creating NewCoupon(): %v", err)
	}

	// Verify code has correct format
	if !strings.HasPrefix(coupon.Code, koText) {
		t.Errorf("Code should start with 'koText'. got: %s", coupon.Code)
	}

	// Check code length
	if len([]rune(coupon.Code)) != maxCodeLength {
		t.Errorf("Code length should be %d. got: %d", maxCodeLength, len([]rune(coupon.Code)))
	}

	// Verify expiration date is set correctly
	if !coupon.ExpireAt.Equal(expiration) {
		t.Errorf("Expiration date not set correctly. expected: %v, got: %v", expiration, coupon.ExpireAt)
	}

	// Check IssuedAt is close to current time
	now := time.Now()
	timeDiff := now.Sub(coupon.IssuedAt)
	if timeDiff > time.Second {
		t.Errorf("IssuedAt time differs too much from current time: %v", timeDiff)
	}
}

func TestCreateCode(t *testing.T) {
	testCases := []struct {
		name    string
		text    string
		nano    int64
		want    string
		wantErr bool
	}{
		{
			name:    "basic test",
			text:    "테스트",
			nano:    1742911351203015000,
			want:    "테스트1203015",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := createCode(tc.text, tc.nano)

			// Check error
			if (err != nil) != tc.wantErr {
				t.Errorf("createCode() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check result
			if got != tc.want {
				t.Errorf("createCode() = %v, want %v", got, tc.want)
			}

			// Check result length
			if len([]rune(got)) > maxCodeLength {
				t.Errorf("Code length is greater than maxCodeLength(%d). Length: %d", maxCodeLength, len([]rune(got)))
			}
		})
	}
}

func TestCreateStringNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		nanoSec  int64
		cnt      int
		expected string
	}{
		{
			name:     "cnt is 5",
			nanoSec:  1234567890000,
			cnt:      5,
			expected: "67890",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := createStringNumbers(tc.nanoSec, tc.cnt)
			if result != tc.expected {
				t.Errorf("createStringNumbers(%d, %d) = %s, want %s",
					tc.nanoSec, tc.cnt, result, tc.expected)
			}
		})
	}
}
