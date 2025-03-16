package cache

import (
	"testing"
	"time"
)

// MockBadgeService implements the BadgeService interface for testing
type MockBadgeService struct {
	badges map[string]*Badge
	calls  int
}

// NewMockBadgeService creates a new MockBadgeService
func NewMockBadgeService() *MockBadgeService {
	return &MockBadgeService{
		badges: make(map[string]*Badge),
		calls:  0,
	}
}

// GetBadge retrieves a badge by ID
func (s *MockBadgeService) GetBadge(id string) (*Badge, error) {
	s.calls++
	return s.badges[id], nil
}

// GetAllBadges retrieves all badges
func (s *MockBadgeService) GetAllBadges() ([]*Badge, error) {
	s.calls++
	badges := make([]*Badge, 0, len(s.badges))
	for _, badge := range s.badges {
		badges = append(badges, badge)
	}
	return badges, nil
}

// CreateBadge creates a new badge
func (s *MockBadgeService) CreateBadge(badge *Badge) error {
	s.calls++
	s.badges[badge.ID] = badge
	return nil
}

// UpdateBadge updates a badge
func (s *MockBadgeService) UpdateBadge(badge *Badge) error {
	s.calls++
	s.badges[badge.ID] = badge
	return nil
}

// DeleteBadge deletes a badge
func (s *MockBadgeService) DeleteBadge(id string) error {
	s.calls++
	delete(s.badges, id)
	return nil
}

// TestBadgeCache_Get tests the Get method of BadgeCache
func TestBadgeCache_Get(t *testing.T) {
	cache := NewBadgeCache(5*time.Minute, 30*time.Second)
	key := "test-key"
	value := "test-value"

	// Test Get when item doesn't exist
	_, found := cache.Get(key)
	if found {
		t.Errorf("Expected item not to be found")
	}

	// Test Get when item exists
	cache.Set(key, value, 0)
	cachedValue, found := cache.Get(key)
	if !found {
		t.Errorf("Expected item to be found")
	}
	if cachedValue != value {
		t.Errorf("Expected %v, got %v", value, cachedValue)
	}

	// Test Get when item has expired
	cache.Set(key, value, 1*time.Nanosecond)
	time.Sleep(10 * time.Millisecond)
	_, found = cache.Get(key)
	if found {
		t.Errorf("Expected item to be expired")
	}
}

// TestBadgeCache_Delete tests the Delete method of BadgeCache
func TestBadgeCache_Delete(t *testing.T) {
	cache := NewBadgeCache(5*time.Minute, 30*time.Second)
	key := "test-key"
	value := "test-value"

	// Test Delete when item exists
	cache.Set(key, value, 0)
	cache.Delete(key)
	_, found := cache.Get(key)
	if found {
		t.Errorf("Expected item to be deleted")
	}
}

// TestCachedBadgeService tests the CachedBadgeService
func TestCachedBadgeService(t *testing.T) {
	mockService := NewMockBadgeService()
	cachedService := NewCachedBadgeService(mockService, 5*time.Minute, 30*time.Second)

	// Create a test badge
	badge := &Badge{
		ID:          "test-badge",
		Name:        "Test Badge",
		Description: "This is a test badge",
		ImageURL:    "https://example.com/badge.png",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test CreateBadge
	err := mockService.CreateBadge(badge)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test GetBadge (first call should hit the service)
	initialCalls := mockService.calls
	retrievedBadge, err := cachedService.GetBadge(badge.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrievedBadge.ID != badge.ID {
		t.Errorf("Expected %v, got %v", badge.ID, retrievedBadge.ID)
	}
	if mockService.calls != initialCalls+1 {
		t.Errorf("Expected service to be called")
	}

	// Test GetBadge (second call should hit the cache)
	initialCalls = mockService.calls
	retrievedBadge, err = cachedService.GetBadge(badge.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrievedBadge.ID != badge.ID {
		t.Errorf("Expected %v, got %v", badge.ID, retrievedBadge.ID)
	}
	if mockService.calls != initialCalls {
		t.Errorf("Expected service not to be called")
	}

	// Test UpdateBadge
	badge.Name = "Updated Test Badge"
	err = cachedService.UpdateBadge(badge)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test GetBadge after update (should hit the cache with updated value)
	retrievedBadge, err = cachedService.GetBadge(badge.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if retrievedBadge.Name != "Updated Test Badge" {
		t.Errorf("Expected %v, got %v", "Updated Test Badge", retrievedBadge.Name)
	}

	// Test DeleteBadge
	err = cachedService.DeleteBadge(badge.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test GetAllBadges (should be empty after deletion)
	badges, err := cachedService.GetAllBadges()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(badges) != 0 {
		t.Errorf("Expected 0 badges, got %v", len(badges))
	}
}
