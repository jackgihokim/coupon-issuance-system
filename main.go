package main

import (
	"github.com/jackgihokim/coupon-issuance-system/server"
)

func main() {
	//const (
	//	campaignName     = "test"
	//	campaignDesc     = "test description"
	//	couponLimitCount = 100000
	//)
	//var (
	//	startDateTime = time.Date(2025, time.March, 26, 0, 0, 0, 0, time.UTC)
	//	endDateTime   = time.Date(2025, time.March, 28, 0, 0, 0, 0, time.UTC)
	//)
	//
	//camp := campaign.NewCampaign(couponLimitCount, campaignName, campaignDesc, startDateTime, endDateTime)
	//fmt.Println(camp.Name)

	srv := server.NewCouponIssuanceServer()
	srv.Start()
}
