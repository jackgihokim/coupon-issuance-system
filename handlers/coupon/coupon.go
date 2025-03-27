package coupon

import (
	"bytes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"

	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
)

const (
	maxCodeLength = 10
	koText        = "테스트"
)

// NewCoupon generates a new Coupon with a unique code, expiration date, and issue timestamp.
// Returns an error if the code generation fails.
func NewCoupon(expiration, now time.Time) (*couponv1.Coupon, error) {
	code, err := createCode(koText, now.UnixNano())
	if err != nil {
		return nil, err
	}

	return &couponv1.Coupon{
		Code:     code,
		ExpireAt: timestamppb.New(expiration),
		IssuedAt: timestamppb.New(now),
	}, nil
}

// createCode generates a string by appending a substring of the nanoseconds to the input text to meet the required length.
// It takes a string `ko` and an int64 `nano` as inputs and returns the generated code or an error if the operation fails.
func createCode(ko string, nano int64) (string, error) {
	r := []rune(ko)
	diff := maxCodeLength - len(r)
	strNum := createStringNumbers(nano, diff)

	var buf bytes.Buffer
	if _, err := buf.WriteString(ko); err != nil {
		return "", err
	}
	if _, err := buf.WriteString(strNum); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// createStringNumbers converts nanoseconds to a microseconds string and extracts the last 'cnt' characters from it.
func createStringNumbers(nanoSec int64, cnt int) string {
	microSec := int(nanoSec) / 1000
	strMicroSec := strconv.Itoa(microSec)
	from := len(strMicroSec) - cnt
	return strMicroSec[from:]
}
