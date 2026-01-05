package memory

import (
	"sync"

	"site-monitor/internal/domain"
)

type StatusMemoryRepository struct {
	mu     sync.RWMutex
	status map[string]domain.SiteCheckStatus
}

func NewStatusMemoryRepository() *StatusMemoryRepository {
	return &StatusMemoryRepository{
		status: make(map[string]domain.SiteCheckStatus),
	}
}
