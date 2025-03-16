package cache

import (
	"sync"
	"time"
)

// CacheItem represents a single item in the cache
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// BadgeCache provides a caching mechanism for badge-related data
type BadgeCache struct {
	items             map[string]CacheItem
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	stopCleanup       chan bool
}

// NewBadgeCache creates a new BadgeCache with the specified default expiration and cleanup interval
func NewBadgeCache(defaultExpiration, cleanupInterval time.Duration) *BadgeCache {
	cache := &BadgeCache{
		items:             make(map[string]CacheItem),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}

	// Start the cleanup routine if cleanup interval is positive
	if cleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// startCleanup periodically removes expired items from the cache
func (c *BadgeCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// deleteExpired removes expired items from the cache
func (c *BadgeCache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration > 0 && v.Expiration < now {
			delete(c.items, k)
		}
	}
}

// Set adds an item to the cache with the specified expiration
func (c *BadgeCache) Set(key string, value interface{}, expiration time.Duration) {
	var exp int64

	if expiration == 0 {
		expiration = c.defaultExpiration
	}

	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: exp,
	}
}

// Get retrieves an item from the cache
func (c *BadgeCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *BadgeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *BadgeCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
}

// Stop stops the cleanup routine
func (c *BadgeCache) Stop() {
	if c.cleanupInterval > 0 {
		c.stopCleanup <- true
	}
}

// Badge represents a badge in the system
type Badge struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	Criteria    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserBadge represents a badge assigned to a user
type UserBadge struct {
	ID        string
	UserID    string
	BadgeID   string
	Badge     *Badge
	AwardedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Event represents a user event in the system
type Event struct {
	ID         string
	UserID     string
	Type       string
	Properties map[string]interface{}
	Timestamp  time.Time
	CreatedAt  time.Time
}

// BadgeProgress represents a user's progress towards earning a badge
type BadgeProgress struct {
	BadgeID      string
	Badge        *Badge
	UserID       string
	Progress     float64
	Requirements map[string]interface{}
	UpdatedAt    time.Time
}

// BadgeService defines the operations for managing badges
type BadgeService interface {
	GetBadge(id string) (*Badge, error)
	GetAllBadges() ([]*Badge, error)
	CreateBadge(badge *Badge) error
	UpdateBadge(badge *Badge) error
	DeleteBadge(id string) error
}

// UserBadgeService defines the operations for managing user badges
type UserBadgeService interface {
	GetUserBadges(userID string) ([]*UserBadge, error)
	AssignBadgeToUser(userBadge *UserBadge) error
}

// EvaluationService defines the operations for badge evaluation
type EvaluationService interface {
	EvaluateUserEvents(userID string, events []*Event) ([]*Badge, error)
	GetUserProgress(userID string) ([]*BadgeProgress, error)
}

// CachedBadgeService provides cached access to badge operations
type CachedBadgeService struct {
	cache       *BadgeCache
	nextService BadgeService
}

// NewCachedBadgeService creates a new CachedBadgeService
func NewCachedBadgeService(nextService BadgeService, defaultExpiration, cleanupInterval time.Duration) *CachedBadgeService {
	return &CachedBadgeService{
		cache:       NewBadgeCache(defaultExpiration, cleanupInterval),
		nextService: nextService,
	}
}

// GetBadge retrieves a badge by ID, using the cache when possible
func (s *CachedBadgeService) GetBadge(id string) (*Badge, error) {
	cacheKey := "badge:" + id

	// Try to get the badge from the cache
	if cachedBadge, found := s.cache.Get(cacheKey); found {
		return cachedBadge.(*Badge), nil
	}

	// Get the badge from the next service
	badge, err := s.nextService.GetBadge(id)
	if err != nil {
		return nil, err
	}

	// Cache the badge
	s.cache.Set(cacheKey, badge, 0) // Use default expiration

	return badge, nil
}

// GetAllBadges retrieves all badges, using the cache when possible
func (s *CachedBadgeService) GetAllBadges() ([]*Badge, error) {
	cacheKey := "badges:all"

	// Try to get the badges from the cache
	if cachedBadges, found := s.cache.Get(cacheKey); found {
		return cachedBadges.([]*Badge), nil
	}

	// Get the badges from the next service
	badges, err := s.nextService.GetAllBadges()
	if err != nil {
		return nil, err
	}

	// Cache the badges
	s.cache.Set(cacheKey, badges, 0) // Use default expiration

	return badges, nil
}

// CreateBadge creates a new badge and invalidates relevant caches
func (s *CachedBadgeService) CreateBadge(badge *Badge) error {
	err := s.nextService.CreateBadge(badge)
	if err != nil {
		return err
	}

	// Invalidate the all badges cache
	s.cache.Delete("badges:all")

	// Cache the new badge
	s.cache.Set("badge:"+badge.ID, badge, 0)

	return nil
}

// UpdateBadge updates a badge and invalidates relevant caches
func (s *CachedBadgeService) UpdateBadge(badge *Badge) error {
	err := s.nextService.UpdateBadge(badge)
	if err != nil {
		return err
	}

	// Invalidate the all badges cache
	s.cache.Delete("badges:all")

	// Update the badge in the cache
	s.cache.Set("badge:"+badge.ID, badge, 0)

	return nil
}

// DeleteBadge deletes a badge and invalidates relevant caches
func (s *CachedBadgeService) DeleteBadge(id string) error {
	err := s.nextService.DeleteBadge(id)
	if err != nil {
		return err
	}

	// Invalidate the all badges cache
	s.cache.Delete("badges:all")

	// Remove the badge from the cache
	s.cache.Delete("badge:" + id)

	return nil
}

// CachedUserBadgeService provides cached access to user badge operations
type CachedUserBadgeService struct {
	cache       *BadgeCache
	nextService UserBadgeService
}

// NewCachedUserBadgeService creates a new CachedUserBadgeService
func NewCachedUserBadgeService(nextService UserBadgeService, defaultExpiration, cleanupInterval time.Duration) *CachedUserBadgeService {
	return &CachedUserBadgeService{
		cache:       NewBadgeCache(defaultExpiration, cleanupInterval),
		nextService: nextService,
	}
}

// GetUserBadges retrieves badges for a user, using the cache when possible
func (s *CachedUserBadgeService) GetUserBadges(userID string) ([]*UserBadge, error) {
	cacheKey := "user_badges:" + userID

	// Try to get the user badges from the cache
	if cachedUserBadges, found := s.cache.Get(cacheKey); found {
		return cachedUserBadges.([]*UserBadge), nil
	}

	// Get the user badges from the next service
	userBadges, err := s.nextService.GetUserBadges(userID)
	if err != nil {
		return nil, err
	}

	// Cache the user badges
	s.cache.Set(cacheKey, userBadges, 0) // Use default expiration

	return userBadges, nil
}

// AssignBadgeToUser assigns a badge to a user and invalidates relevant caches
func (s *CachedUserBadgeService) AssignBadgeToUser(userBadge *UserBadge) error {
	err := s.nextService.AssignBadgeToUser(userBadge)
	if err != nil {
		return err
	}

	// Invalidate the user badges cache
	s.cache.Delete("user_badges:" + userBadge.UserID)

	return nil
}

// CachedEvaluationService provides cached access to badge evaluation operations
type CachedEvaluationService struct {
	cache       *BadgeCache
	nextService EvaluationService
}

// NewCachedEvaluationService creates a new CachedEvaluationService
func NewCachedEvaluationService(nextService EvaluationService, defaultExpiration, cleanupInterval time.Duration) *CachedEvaluationService {
	return &CachedEvaluationService{
		cache:       NewBadgeCache(defaultExpiration, cleanupInterval),
		nextService: nextService,
	}
}

// EvaluateUserEvents evaluates user events for badge eligibility, using the cache when possible
func (s *CachedEvaluationService) EvaluateUserEvents(userID string, events []*Event) ([]*Badge, error) {
	// For evaluation, we typically don't want to cache the results as events change frequently
	// However, we can cache badge criteria and rules

	return s.nextService.EvaluateUserEvents(userID, events)
}

// GetUserProgress retrieves a user's progress towards badges, using the cache when possible
func (s *CachedEvaluationService) GetUserProgress(userID string) ([]*BadgeProgress, error) {
	cacheKey := "user_progress:" + userID

	// Try to get the user progress from the cache
	if cachedProgress, found := s.cache.Get(cacheKey); found {
		return cachedProgress.([]*BadgeProgress), nil
	}

	// Get the user progress from the next service
	progress, err := s.nextService.GetUserProgress(userID)
	if err != nil {
		return nil, err
	}

	// Cache the user progress
	s.cache.Set(cacheKey, progress, 0) // Use default expiration

	return progress, nil
}

// InvalidateUserCache invalidates all cache entries for a user
func (s *CachedEvaluationService) InvalidateUserCache(userID string) {
	s.cache.Delete("user_progress:" + userID)
}
