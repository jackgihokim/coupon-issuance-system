package server

import (
	"context"
	"errors"
	"github.com/jackgihokim/coupon-issuance-system/common/id"
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/jackgihokim/coupon-issuance-system/handlers/campaign"
	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
	"github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1/couponv1connect"
)

type CouponIssuanceServer struct{}

const httpAddr = "localhost:8080"

// NewCouponIssuanceServer initializes and returns a new instance of CouponIssuanceServer.
func NewCouponIssuanceServer() *CouponIssuanceServer {
	return &CouponIssuanceServer{}
}

// Start initializes the HTTP server, sets up routes for the CouponIssuanceService, and begins listening for requests.
func (s *CouponIssuanceServer) Start() {
	mux := http.NewServeMux()
	path, handler := couponv1connect.NewCouponIssuanceServiceHandler(s)
	mux.Handle(path, handler)
	err := http.ListenAndServe(httpAddr, h2c.NewHandler(mux, &http2.Server{}))
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	campaignId *id.ID
	coupons    *coupon.Coupons
)

// CreateCampaign creates a new campaign with the specified details and returns the created campaign or an error.
func (s *CouponIssuanceServer) CreateCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.CreateCampaignRequest],
) (*connect.Response[couponv1.CreateCampaignResponse], error) {
	campaignId = id.NewID()
	coupons = coupon.NewCoupons(req.Msg.CouponLimit)

	camp, err := campaign.NewCampaign(
		campaignId, req.Msg.CouponLimit, req.Msg.Name, req.Msg.Description, req.Msg.StartAt, req.Msg.EndAt, coupons,
	)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&couponv1.CreateCampaignResponse{
		Campaign: camp,
	})
	return resp, nil
}

// GetCampaign retrieves a campaign by its unique ID using the provided request and returns the campaign details or an error.
func (s *CouponIssuanceServer) GetCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.GetCampaignRequest],
) (*connect.Response[couponv1.GetCampaignResponse], error) {
	camp, err := campaign.GetCampaign(req.Msg.CampaignId)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&couponv1.GetCampaignResponse{
		Campaign: camp,
	})
	return resp, nil
}

// IssueCoupon issues a new coupon for the specified campaign if the campaign is active and within its validity period.
func (s *CouponIssuanceServer) IssueCoupon(
	ctx context.Context,
	req *connect.Request[couponv1.IssueCouponRequest],
) (*connect.Response[couponv1.IssueCouponResponse], error) {
	camp, err := campaign.GetCampaign(req.Msg.CampaignId)
	err = validatePeriod(camp)
	if err != nil {
		return nil, err
	}

	coup, err := coupon.NewCoupon()
	if err != nil {
		return nil, err
	}

	// TODO: needs to push coupon to coupons list.

	resp := connect.NewResponse(&couponv1.IssueCouponResponse{
		Coupon: coup,
	})
	return resp, nil
}

// validatePeriod validates if the campaign is within its validity period by checking the start and end times.
// Returns an error if the campaign has not started or has already ended.
func validatePeriod(camp *couponv1.Campaign) error {
	if camp.StartAt.AsTime().Before(time.Now()) {
		return errors.New("campaign is not started yet")
	}
	if camp.EndAt.AsTime().Before(time.Now()) {
		return errors.New("campaign is over")
	}
	return nil
}
