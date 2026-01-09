package memory

import (
	"sync"

	"site-monitor/internal/domain"
)

type StatusMemoryRepository struct {
	mu     sync.RWMutex
	status map[string]domain.SiteCheckStatus
}

func (r *StatusMemoryRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = make(map[string]domain.SiteCheckStatus)
}

func (r *StatusMemoryRepository) Init(statuses []domain.SiteCheckStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.status = make(map[string]domain.SiteCheckStatus)
	for _, s := range statuses {
		r.status[s.SiteID] = s
	}
}

func NewStatusMemoryRepository() *StatusMemoryRepository {
	return &StatusMemoryRepository{
		status: make(map[string]domain.SiteCheckStatus),
	}
}

func (r *StatusMemoryRepository) Save(s domain.SiteCheckStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status[s.SiteID] = s
}

func (r *StatusMemoryRepository) GetBySiteID(siteID string) (*domain.SiteCheckStatus, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.status[siteID]
	if !ok {
		return nil, false
	}

	return &s, true
}