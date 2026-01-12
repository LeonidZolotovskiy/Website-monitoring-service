package memory

import (
	"sync"

	"site-monitor/internal/domain"
	"site-monitor/internal/repository"
)

type SiteMemoryRepository struct {
	mu    sync.RWMutex
	sites map[string]domain.Site
	sitesByURL map[string]string 
}


func (r *SiteMemoryRepository) Reset() {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.sites = make(map[string]domain.Site)
}

func PopulateRepository(repo *SiteMemoryRepository, sites []domain.Site) error {
    for _, s := range sites {
        if err := repo.Create(s); err != nil {
            return err
        }
    }
    return nil
}

func (r *SiteMemoryRepository) Init(sites []domain.Site) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sites = make(map[string]domain.Site)
	for _, s := range sites {
		r.sites[s.ID] = s
	}
}

func NewSiteMemoryRepository() *SiteMemoryRepository {
	return &SiteMemoryRepository{
		sites: make(map[string]domain.Site),
		sitesByURL: make(map[string]string),
	}
}

func (r *SiteMemoryRepository) Create(site domain.Site) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.sitesByURL[site.URL]; exists {
        return repository.ErrSiteAlreadyExists
    }

    r.sites[site.ID] = site
    r.sitesByURL[site.URL] = site.ID

    return nil
}

func (r *SiteMemoryRepository) GetByURL(url string) (*domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	site, exists := r.sites[url]
	if !exists {
		return nil, repository.ErrSiteNotFound
	}

	return &site, nil
}

func (r *SiteMemoryRepository) GetAll() ([]domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Site, 0, len(r.sites))
	for _, site := range r.sites {
		result = append(result, site)
	}

	return result, nil
}

func (r *SiteMemoryRepository) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	site, ok := r.sites[id]
	if !ok {
		return repository.ErrSiteNotFound
	}

	delete(r.sites, id)
	delete(r.sitesByURL, site.URL)

	return nil
}


func (r *SiteMemoryRepository) GetByID(id string) (*domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	site, ok := r.sites[id]
	if !ok {
		return nil, repository.ErrSiteNotFound
	}

	return &site, nil
}