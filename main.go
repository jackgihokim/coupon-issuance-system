package main

import (
	"github.com/jackgihokim/coupon-issuance-system/server"
)

func main() {
	srv := server.NewCouponIssuanceServer()
	srv.Start()
}
