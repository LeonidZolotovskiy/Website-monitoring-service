package memory

import (
	"sync"

	"site-monitor/internal/domain"
)

type SiteMemoryRepository struct {
	mu    sync.RWMutex
	sites []domain.Site
}

func NewSiteMemoryRepository(initial []domain.Site) *SiteMemoryRepository {
	return &SiteMemoryRepository{
		sites: initial,
	}
}

func (r *SiteMemoryRepository) GetAll() []domain.Site {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// важно: возвращаем копию, а не оригинальный slice
	result := make([]domain.Site, len(r.sites))
	copy(result, r.sites)

	return result
}
