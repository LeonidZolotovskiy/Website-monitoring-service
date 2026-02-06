package memory

import (
	"sync"
	"context"
	"errors"
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

func (r *StatusMemoryRepository) GetBySiteID(ctx context.Context,siteID string) (*domain.SiteCheckStatus, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.status[siteID]
	if !ok {
		return nil, false
	}

	return &s, true
}

func (r *StatusMemoryRepository) GetHistoryBySiteID(ctx context.Context,siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    s, ok := r.status[siteID]
    if !ok {
        return nil, 0, errors.New("site not found")
    }

    total := 1

    if offset >= total {
        return []domain.SiteCheckStatus{}, total, nil
    }

    history := []domain.SiteCheckStatus{s}

    return history, total, nil
}
