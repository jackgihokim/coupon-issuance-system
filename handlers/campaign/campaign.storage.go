package campaign

import (
	"errors"
	"sync"
)

// test-store for the created campaigns.
type Store struct {
	mu sync.Mutex
	m  map[uint32]*Campaign
}

// newCampaignStore initializes and returns a new instance of Store with an empty campaign map.
func newCampaignStore() *Store {
	return &Store{
		m: make(map[uint32]*Campaign),
	}
}

// add adds a new campaign to the store. It safely locks the store to ensure thread-safe operations.
func (s *Store) add(campaign *Campaign) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[campaign.Id] = campaign
	return nil
}

// delete removes the campaign with the specified ID from the store in a thread-safe manner.
// Returns an error if any issues occur.
func (s *Store) delete(id uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, id)
	return nil
}

// get retrieves a campaign by its ID from the store in a thread-safe manner.
// Returns an error if the campaign is not found.
func (s *Store) get(id uint32) (*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	camp, ok := s.m[id]
	if !ok {
		return nil, errors.New("campaign not found")
	}
	return camp, nil
}

// list returns a thread-safe list of all campaigns stored in the store.
func (s *Store) list() ([]*Campaign, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []*Campaign
	for _, v := range s.m {
		list = append(list, v)
	}
	return list, nil
}
