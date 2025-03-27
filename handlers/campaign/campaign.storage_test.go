package campaign

import (
	"testing"
)

func TestNewCampaignStore(t *testing.T) {
	store := newCampaignStore()
	if store == nil {
		t.Errorf("newCampaignStore() returned nil")
	}
	if store.m == nil {
		t.Errorf("initialized store map is nil")
	}
	if len(store.m) != 0 {
		t.Errorf("initialized store map is not empty, count: %d", len(store.m))
	}
}

func TestAdd(t *testing.T) {
	store := newCampaignStore()
	campaign := &Campaign{Id: 1}

	err := store.add(campaign)
	if err != nil {
		t.Errorf("error occurred while adding campaign: %v", err)
	}

	// Verify campaign was added to the map
	if len(store.m) != 1 {
		t.Errorf("campaign was not properly added, current map size: %d", len(store.m))
	}

	storedCampaign, exists := store.m[1]
	if !exists {
		t.Errorf("campaign with ID 1 does not exist in map")
	}

	if storedCampaign != campaign {
		t.Errorf("stored campaign doesn't match the original")
	}
}

func TestGet(t *testing.T) {
	store := newCampaignStore()
	campaign := &Campaign{Id: 1}
	store.m[campaign.Id] = campaign

	// Getting an existing campaign
	gotCampaign, err := store.get(1)
	if err != nil {
		t.Errorf("error occurred while getting existing campaign: %v", err)
	}
	if gotCampaign != campaign {
		t.Errorf("returned campaign doesn't match the original")
	}

	// Getting a non-existent campaign
	_, err = store.get(2)
	if err == nil {
		t.Errorf("expected error when getting non-existent campaign, but got nil")
	}
}

func TestDelete(t *testing.T) {
	store := newCampaignStore()
	campaign := &Campaign{Id: 1}
	store.m[campaign.Id] = campaign

	err := store.delete(1)
	if err != nil {
		t.Errorf("error occurred while deleting campaign: %v", err)
	}

	// Verify deletion
	if _, exists := store.m[1]; exists {
		t.Errorf("campaign still exists after deletion")
	}
}

func TestList(t *testing.T) {
	store := newCampaignStore()

	// Test with empty store
	campaigns, err := store.list()
	if err != nil {
		t.Errorf("error occurred while listing from empty store: %v", err)
	}
	if len(campaigns) != 0 {
		t.Errorf("list from empty store is not empty, count: %d", len(campaigns))
	}

	// Add campaigns
	campaign1 := &Campaign{Id: 1}
	campaign2 := &Campaign{Id: 2}
	store.m[campaign1.Id] = campaign1
	store.m[campaign2.Id] = campaign2

	campaigns, err = store.list()
	if err != nil {
		t.Errorf("error occurred while listing campaigns: %v", err)
	}
	if len(campaigns) != 2 {
		t.Errorf("expected campaign count: 2, actual: %d", len(campaigns))
	}

	// Verify all campaigns are included in the list
	campaignMap := make(map[uint32]*Campaign)
	for _, c := range campaigns {
		campaignMap[c.Id] = c
	}

	if _, exists := campaignMap[1]; !exists {
		t.Errorf("campaign with ID 1 is missing from list")
	}

	if _, exists := campaignMap[2]; !exists {
		t.Errorf("campaign with ID 2 is missing from list")
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := newCampaignStore()
	done := make(chan bool)

	// Access from multiple goroutines concurrently
	for i := 0; i < 10; i++ {
		id := uint32(i)
		go func(id uint32) {
			campaign := &Campaign{Id: id}
			store.add(campaign)
			done <- true
		}(id)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all campaigns were added correctly
	campaigns, _ := store.list()
	if len(campaigns) != 10 {
		t.Errorf("expected campaign count: 10, actual: %d", len(campaigns))
	}
}

func TestAddAndGet(t *testing.T) {
	store := newCampaignStore()
	campaign := &Campaign{Id: 42}

	// Add campaign
	err := store.add(campaign)
	if err != nil {
		t.Errorf("failed to add campaign: %v", err)
	}

	// Retrieve and verify
	retrieved, err := store.get(42)
	if err != nil {
		t.Errorf("failed to get campaign: %v", err)
	}

	if retrieved != campaign {
		t.Errorf("retrieved campaign doesn't match the added campaign")
	}
}

func TestDeleteAndGet(t *testing.T) {
	store := newCampaignStore()
	campaign := &Campaign{Id: 42}

	// Add and then delete
	store.add(campaign)
	err := store.delete(42)
	if err != nil {
		t.Errorf("failed to delete campaign: %v", err)
	}

	// Try to get deleted campaign
	_, err = store.get(42)
	if err == nil {
		t.Errorf("expected error when getting deleted campaign")
	}
}
