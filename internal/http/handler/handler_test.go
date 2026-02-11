package handler_test

import (
	"time"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"site-monitor/internal/domain"
	"site-monitor/internal/http/handler"
	"github.com/go-chi/chi/v5"
	"site-monitor/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"
)

//
// ---------- MOCK REPOSITORY ----------
//

type MockSiteRepository struct {
	CreateFunc     func(ctx context.Context, site domain.Site) (string, error)
	GetByURLFunc   func(ctx context.Context, url string) (*domain.Site, error)
	GetAllFunc     func(ctx context.Context) ([]domain.Site, error)
	DeleteByIDFunc func(ctx context.Context, logger *slog.Logger, id string) error
	GetByIDFunc    func(ctx context.Context, id string) (*domain.Site, error)
}

func (m *MockSiteRepository) Create(ctx context.Context, site domain.Site) (string, error) {
	return m.CreateFunc(ctx, site)
}

func (m *MockSiteRepository) GetByURL(ctx context.Context, url string) (*domain.Site, error) {
	return m.GetByURLFunc(ctx, url)
}

func (m *MockSiteRepository) GetAll(ctx context.Context) ([]domain.Site, error) {
	return m.GetAllFunc(ctx)
}

func (m *MockSiteRepository) DeleteByID(ctx context.Context, logger *slog.Logger, id string) error {
	return m.DeleteByIDFunc(ctx, logger, id)
}

func (m *MockSiteRepository) GetByID(ctx context.Context, id string) (*domain.Site, error) {
	return m.GetByIDFunc(ctx, id)
}

// ---------- MOCK STATUS REPOSITORY ----------
type MockStatusRepository struct {
	GetBySiteIDFunc       func(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool)
	GetHistoryBySiteIDFunc func(ctx context.Context, siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error)
	SaveFunc              func(status domain.SiteCheckStatus)
}

func (m *MockStatusRepository) GetBySiteID(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool) {
	return m.GetBySiteIDFunc(ctx, siteID)
}

func (m *MockStatusRepository) GetHistoryBySiteID(ctx context.Context, siteID string, limit, offset int) ([]domain.SiteCheckStatus, int, error) {
	return m.GetHistoryBySiteIDFunc(ctx, siteID, limit, offset)
}

func (m *MockStatusRepository) Save(status domain.SiteCheckStatus) {
	if m.SaveFunc != nil {
		m.SaveFunc(status)
	}
}
//
// ---------- ROUTER HELPER ----------
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
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func setupCreateHandler(mockRepo *MockSiteRepository) *handler.SiteHandler {
	return handler.NewSiteHandler(mockRepo, nil, nil, newLogger())
}
//
// ---------- TESTS ----------
//

func TestGetSitesHandler(t *testing.T) {
	mockRepo := &MockSiteRepository{
		GetAllFunc: func(ctx context.Context) ([]domain.Site, error) {
			return []domain.Site{
				{ID: "1", Name: "Google", URL: "https://google.com"},
				{ID: "2", Name: "GitHub", URL: "https://github.com"},
			}, nil
		},
	}

	h := setupCreateHandler(mockRepo)
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

//
// ---------- CREATE ----------
//

func TestCreateSiteHandler(t *testing.T) {
	mockRepo := &MockSiteRepository{
		CreateFunc: func(ctx context.Context, site domain.Site) (string, error) {
			// Возвращаем ID, как будто сайт успешно создан
			return uuid.NewString(), nil
		},
	}

	h := setupCreateHandler(mockRepo)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "success",
			body:       `{"name":"StackOverflow","url":"https://stackoverflow.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "empty URL",
			body:       `{"name":"Example","url":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid URL format",
			body:       `{"name":"Example","url":"not-a-valid-url"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty JSON",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sites", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.Create(w, req)

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

// //
// // ---------- DELETE ----------
// //

func TestDeleteSiteHandler(t *testing.T) {
	mockRepo := &MockSiteRepository{
		DeleteByIDFunc: func(ctx context.Context, logger *slog.Logger, id string) error {
			if id == "1" {
				return nil 
			}
			return repository.ErrSiteNotFound 
		},
	}

	h := setupCreateHandler(mockRepo)
	mux := setupRouter(h)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "success",
			id:         "1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "not found",
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


//
// ---------- STATUS ----------
//

func TestGetSiteStatusHandler(t *testing.T) {
    mockRepo := &MockSiteRepository{
        GetByIDFunc: func(ctx context.Context, id string) (*domain.Site, error) {
            if id == "1" {
                return &domain.Site{
                    ID:   "1",
                    Name: "Google",
                    URL:  "https://google.com",
                }, nil
            }
            return nil, errors.New("not found")
        },
    }

    mockStatusRepo := &MockStatusRepository{
        GetBySiteIDFunc: func(ctx context.Context, siteID string) (*domain.SiteCheckStatus, bool) {
            if siteID == "1" {
                rt := time.Duration(123) * time.Millisecond
                lastChecked := time.Now()
                code := 200
                return &domain.SiteCheckStatus{
                    URL:           "https://google.com",
                    Status:        domain.StatusOK,
                    StatusCode:    &code,
                    ResponseTime:  &rt,
                    LastCheckedAt: &lastChecked,
                    Error:         "",
                }, true
            }
            return nil, false
        },
    }

    h := handler.NewSiteHandler(mockRepo, mockStatusRepo, nil, newLogger())

    // chi роутер
    r := chi.NewRouter()
    r.Get("/api/v1/sites/{id}/status", h.GetStatus)

    // тестовый запрос
    req := httptest.NewRequest("GET", "/api/v1/sites/1/status", nil)
    w := httptest.NewRecorder()

    r.ServeHTTP(w, req)

    resp := w.Result()
    defer resp.Body.Close()

    assert.Equal(t, 200, resp.StatusCode)

    var statusResp handler.SiteStatusResponse
    err := json.NewDecoder(resp.Body).Decode(&statusResp)
    assert.NoError(t, err)
    assert.Equal(t, "https://google.com", statusResp.URL)
    assert.Equal(t, string(domain.StatusOK), statusResp.Status)
    assert.NotNil(t, statusResp.StatusCode)
    assert.NotNil(t, statusResp.ResponseTime)
    assert.NotNil(t, statusResp.LastCheckedAt)
}



