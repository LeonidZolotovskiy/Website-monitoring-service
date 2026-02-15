package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"site-monitor/internal/domain"
	"site-monitor/internal/http/handler"
	"site-monitor/internal/repository"
)

//
// ---------- MOCK REPOSITORIES ----------
//

type MockSiteRepository struct {
	CreateFunc     func(ctx context.Context, site domain.Site) (string, error)
	GetByURLFunc   func(ctx context.Context, url string) (*domain.Site, error)
	GetAllFunc     func(ctx context.Context) ([]domain.Site, error)
	DeleteByIDFunc func(ctx context.Context, logger *slog.Logger, id string) error
	GetByIDFunc    func(ctx context.Context, id string) (*domain.Site, error)
}

func (m *MockSiteRepository) Create(ctx context.Context, site domain.Site) (string, error) {
	if m.CreateFunc == nil {
		return "", errors.New("CreateFunc is not set")
	}
	return m.CreateFunc(ctx, site)
}

func (m *MockSiteRepository) GetByURL(ctx context.Context, url string) (*domain.Site, error) {
	if m.GetByURLFunc == nil {
		return nil, errors.New("GetByURLFunc is not set")
	}
	return m.GetByURLFunc(ctx, url)
}

func (m *MockSiteRepository) GetAll(ctx context.Context) ([]domain.Site, error) {
	if m.GetAllFunc == nil {
		return nil, errors.New("GetAllFunc is not set")
	}
	return m.GetAllFunc(ctx)
}

func (m *MockSiteRepository) DeleteByID(ctx context.Context, logger *slog.Logger, id string) error {
	if m.DeleteByIDFunc == nil {
		return errors.New("DeleteByIDFunc is not set")
	}
	return m.DeleteByIDFunc(ctx, logger, id)
}

func (m *MockSiteRepository) GetByID(ctx context.Context, id string) (*domain.Site, error) {
	if m.GetByIDFunc == nil {
		return nil, errors.New("GetByIDFunc is not set")
	}
	return m.GetByIDFunc(ctx, id)
}

type MockStatusRepository struct {
	GetBySiteIDFunc        func(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool)
	GetHistoryBySiteIDFunc func(ctx context.Context, siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error)
	SaveFunc               func(status domain.SiteCheckStatus)
}

func (m *MockStatusRepository) GetBySiteID(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool) {
	if m.GetBySiteIDFunc == nil {
		return nil, false
	}
	return m.GetBySiteIDFunc(ctx, siteID)
}

func (m *MockStatusRepository) GetHistoryBySiteID(ctx context.Context, siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error) {
	if m.GetHistoryBySiteIDFunc == nil {
		return nil, 0, errors.New("GetHistoryBySiteIDFunc is not set")
	}
	return m.GetHistoryBySiteIDFunc(ctx, siteID, limit, offset)
}

func (m *MockStatusRepository) Save(status domain.SiteCheckStatus) {
	if m.SaveFunc != nil {
		m.SaveFunc(status)
	}
}

//
// ---------- TEST HELPERS ----------
//

func setupRouter(h *handler.SiteHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/sites", h.GetSites)
	mux.HandleFunc("POST /api/v1/sites", h.Create)
	mux.HandleFunc("DELETE /api/v1/sites/{id}", h.Delete)
	mux.HandleFunc("GET /api/v1/sites/{id}/status", h.GetStatus)

	return mux
}

func newLogger() *slog.Logger {
	// Чтобы тесты не шумели в stdout
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newHandler(siteRepo *MockSiteRepository, statusRepo *MockStatusRepository) *handler.SiteHandler {
	return handler.NewSiteHandler(siteRepo, statusRepo, nil, newLogger())
}

//
// ---------- TESTS ----------
//

func TestGetSites_Success(t *testing.T) {
	mockRepo := &MockSiteRepository{
		GetAllFunc: func(ctx context.Context) ([]domain.Site, error) {
			return []domain.Site{
				{ID: "1", Name: "Google", URL: "https://google.com"},
				{ID: "2", Name: "GitHub", URL: "https://github.com"},
			}, nil
		},
	}

	h := newHandler(mockRepo, nil)
	mux := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sites", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var sites []domain.Site
	err := json.NewDecoder(resp.Body).Decode(&sites)

	assert.NoError(t, err)
	assert.Len(t, sites, 2)
	assert.Equal(t, "Google", sites[0].Name)
}

func TestCreateSite(t *testing.T) {
	mockRepo := &MockSiteRepository{
		CreateFunc: func(ctx context.Context, site domain.Site) (string, error) {
			return uuid.NewString(), nil
		},
	}

	h := newHandler(mockRepo, nil)
	mux := setupRouter(h)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Success",
			body:       `{"name":"StackOverflow","url":"https://stackoverflow.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "EmptyURL_BadRequest",
			body:       `{"name":"Example","url":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidURLFormat_BadRequest",
			body:       `{"name":"Example","url":"not-a-valid-url"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "EmptyJSON_BadRequest",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sites", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusCreated {
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

				var site domain.Site
				err := json.NewDecoder(resp.Body).Decode(&site)
				assert.NoError(t, err)
				assert.NotEmpty(t, site.ID)
			}
		})
	}
}

func TestDeleteSite(t *testing.T) {
	mockRepo := &MockSiteRepository{
		DeleteByIDFunc: func(ctx context.Context, logger *slog.Logger, id string) error {
			if id == "1" {
				return nil
			}
			return repository.ErrSiteNotFound
		},
	}

	h := newHandler(mockRepo, nil)
	mux := setupRouter(h)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "Success_NoContent",
			id:         "1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "NotFound",
			id:         "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/sites/"+tt.id, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Result().StatusCode)
		})
	}
}

func TestGetSiteStatus_Success(t *testing.T) {
	mockRepo := &MockSiteRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*domain.Site, error) {
			if id == "1" {
				return &domain.Site{
					ID:   "1",
					Name: "Google",
					URL:  "https://google.com",
				}, nil
			}
			return nil, repository.ErrSiteNotFound
		},
	}

	fixedCheckedAt := time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)

	mockStatusRepo := &MockStatusRepository{
		GetBySiteIDFunc: func(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool) {
			if siteID == "1" {
				rt := 123 * time.Millisecond
				code := 200
				return &domain.SiteCheckStatus{
					URL:           "https://google.com",
					Status:        domain.StatusOK,
					StatusCode:    &code,
					ResponseTime:  &rt,
					LastCheckedAt: &fixedCheckedAt,
					Error:         "",
				}, true
			}
			return nil, false
		},
	}

	h := newHandler(mockRepo, mockStatusRepo)
	mux := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sites/1/status", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var got handler.SiteStatusResponse
	err := json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)

	// обязательные/ожидаемые в success
	assert.Equal(t, "https://google.com", got.URL)
	assert.Equal(t, "ok", got.Status)

	// lastCheckedAt и responseTimeMs должны быть
	if assert.NotNil(t, got.LastCheckedAt) {
		assert.True(t, got.LastCheckedAt.Equal(fixedCheckedAt))
	}

	if assert.NotNil(t, got.ResponseTime) {
		assert.EqualValues(t, 123, *got.ResponseTime)
	}

	// error в success быть не должно (omitempty => nil)
	assert.Nil(t, got.Error)

	// statusCode — опционально (у тебя в реальном RAW его не было)
	// Если хочешь — можешь сделать "если есть, то 200"
	if got.StatusCode != nil {
		assert.Equal(t, 200, *got.StatusCode)
	}
}




