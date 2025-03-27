package campaign

import (
	"time"

	"github.com/jackgihokim/coupon-issuance-system/common/id"
	"github.com/jackgihokim/coupon-issuance-system/handlers/coupon"
)

type Campaign struct {
	Id          uint32
	CouponLimit uint32
	Name        string
	Description string
	CreatedAt   time.Time
	StartAt     time.Time
	EndAt       time.Time
	Coupons     *coupon.Coupons
}

var CampaignId *id.ID
var store *Store = newCampaignStore()

// NewCampaign creates a new campaign with the provided parameters and stores it.
// Returns a pointer to the newly created Campaign object or an error if the campaign could not be stored.
func NewCampaign(limit uint32, name, desc string, start, end time.Time) (*Campaign, error) {
	camp := &Campaign{
		Id:          CampaignId.Next(),
		CouponLimit: limit,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now().UTC(), // must use UTC for being the same as timestamppb.
		StartAt:     start,
		EndAt:       end,
		Coupons:     coupon.NewCoupons(limit),
	}

	err := store.add(camp)
	if err != nil {
		return nil, err
	}

	return camp, nil
}

// GetCampaign retrieves a campaign by its unique ID from the store.
// Returns the campaign details or an error if the campaign does not exist.
func GetCampaign(id uint32) (*Campaign, error) {
	camp, err := store.get(id)
	if err != nil {
		return nil, err
	}
	return camp, nil
}
