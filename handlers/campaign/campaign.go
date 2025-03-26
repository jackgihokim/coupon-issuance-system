package campaign

import (
	"github.com/jackgihokim/coupon-issuance-system/common/id"
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
	couponv1 "github.com/jackgihokim/coupon-issuance-system/protos/coupon/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var storage *Storage = newCampaignStorage()

// NewCampaign creates a new campaign with the provided details and stores it. Returns the created campaign or an error.
// id is a pointer to the ID generator for unique campaign ID creation.
// limit specifies the maximum number of coupons allowed in the campaign.
// name and desc are strings representing the campaign's name and description, respectively.
// start and end are pointers to timestamps defining the campaign's start and end time.
// coups is a pointer to a Coupons list containing coupon details for the campaign.
// Returns a pointer to the created couponv1.Campaign or an error if the storage operation fails.
func NewCampaign(
	id *id.ID, limit uint64, name, desc string, start, end *timestamppb.Timestamp, coups *coupon.Coupons,
) (*couponv1.Campaign, error) {
	camp := &couponv1.Campaign{
		Id:          id.Next(),
		CouponLimit: limit,
		Name:        name,
		Description: desc,
		CreatedAt:   timestamppb.Now(),
		StartAt:     start,
		EndAt:       end,
		Coupons:     coups.List(),
	}

	err := storage.add(camp)
	if err != nil {
		return nil, err
	}

	return camp, nil
}

// GetCampaign retrieves a campaign by its unique ID from the storage.
// It returns the campaign object and an error if the operation fails.
func GetCampaign(id uint64) (*couponv1.Campaign, error) {
	camp, err := storage.get(id)
	if err != nil {
		return nil, err
	}
	return camp, nil
}
