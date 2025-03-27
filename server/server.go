package server

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/jackgihokim/coupon-issuance-system/handlers/campaign"
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
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

// CreateCampaign handles the creation of a new campaign with provided details.
// Returns the created campaign or an error.
func (s *CouponIssuanceServer) CreateCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.CreateCampaignRequest],
) (*connect.Response[couponv1.CreateCampaignResponse], error) {
	camp, err := campaign.NewCampaign(
		req.Msg.CouponLimit, req.Msg.Name, req.Msg.Description, req.Msg.StartAt.AsTime(), req.Msg.EndAt.AsTime(),
	)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&couponv1.CreateCampaignResponse{
		Campaign: &couponv1.Campaign{
			Id:          camp.Id,
			CouponLimit: camp.CouponLimit,
			Name:        camp.Name,
			Description: camp.Description,
			CreatedAt:   timestamppb.New(camp.CreatedAt),
			StartAt:     timestamppb.New(camp.StartAt),
			EndAt:       timestamppb.New(camp.EndAt),
			Coupons:     camp.Coupons.List(),
		},
	})
	return resp, nil
}

// GetCampaign retrieves the details of a specific campaign using the provided campaign ID.
// Returns a response containing the campaign details or an error if the campaign is not found.
func (s *CouponIssuanceServer) GetCampaign(
	ctx context.Context,
	req *connect.Request[couponv1.GetCampaignRequest],
) (*connect.Response[couponv1.GetCampaignResponse], error) {
	camp, err := campaign.GetCampaign(req.Msg.CampaignId)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&couponv1.GetCampaignResponse{
		Campaign: &couponv1.Campaign{
			Id:          camp.Id,
			CouponLimit: camp.CouponLimit,
			Name:        camp.Name,
			Description: camp.Description,
			CreatedAt:   timestamppb.New(camp.CreatedAt),
			StartAt:     timestamppb.New(camp.StartAt),
			EndAt:       timestamppb.New(camp.EndAt),
			Coupons:     camp.Coupons.List(),
		},
	})
	return resp, nil
}

// IssueCoupon handles the issuance of a new coupon for a specific campaign, validating campaign status and period.
// Returns a response containing the issued coupon or an error if the operation fails.
func (s *CouponIssuanceServer) IssueCoupon(
	ctx context.Context,
	req *connect.Request[couponv1.IssueCouponRequest],
) (*connect.Response[couponv1.IssueCouponResponse], error) {
	camp, err := campaign.GetCampaign(req.Msg.CampaignId)
	now := time.Now().UTC() // must use UTC for being the same as timestamppb.
	err = validatePeriod(camp, now)
	if err != nil {
		return nil, err
	}

	coup, err := coupon.NewCoupon(now)
	if err != nil {
		return nil, err
	}

	err = camp.Coupons.Add(coup)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&couponv1.IssueCouponResponse{
		Coupon: coup,
	})
	return resp, nil
}

// validatePeriod checks if the given campaign is currently active based on its start and end times relative to now.
// Returns an error if the campaign has not started or has already ended.
func validatePeriod(camp *campaign.Campaign, now time.Time) error {
	if camp.StartAt.After(now) {
		return errors.New("campaign is not started yet")
	}
	if camp.EndAt.Before(now) {
		return errors.New("campaign is over")
	}
	return nil
}
