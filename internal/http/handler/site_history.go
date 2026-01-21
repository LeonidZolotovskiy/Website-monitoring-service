package handler

import (
    "net/http"
    "strconv"
    "site-monitor/internal/repository"
    "log/slog"

    "github.com/go-chi/chi/v5"
)

const (
    DefaultLimit = 20
    MaxLimit     = 100
)

func (h *SiteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    siteID := chi.URLParam(r, "id")
    if siteID == "" {
        http.Error(w, "missing site id", http.StatusBadRequest)
        return
    }

    _, err := h.siteRepo.GetByID(r.Context(), siteID)
    if err != nil {
        if err == repository.ErrSiteNotFound {
            http.Error(w, "site not found", http.StatusNotFound)
            return
        }
        h.logger.Error("Failed to get site", slog.String("siteID", siteID), slog.Any("err", err))
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    query := r.URL.Query()
    limit := DefaultLimit
    offset := 0

    if l := query.Get("limit"); l != "" {
        if val, err := strconv.Atoi(l); err == nil && val > 0 {
            if val > MaxLimit {
                limit = MaxLimit
            } else {
                limit = val
            }
        } else {
            http.Error(w, "invalid limit parameter", http.StatusBadRequest)
            return
        }
    }

    if o := query.Get("offset"); o != "" {
        if val, err := strconv.Atoi(o); err == nil && val >= 0 {
            offset = val
        } else {
            http.Error(w, "invalid offset parameter", http.StatusBadRequest)
            return
        }
    }

    history, total, err := h.statusRepo.GetHistoryBySiteID(siteID, limit, offset)
    if err != nil {
        h.logger.Error("Failed to get site history", slog.String("siteID", siteID), slog.Any("err", err))
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    var items []SiteCheckHistoryItem
    for _, h := range history {
        items = append(items, SiteCheckHistoryItem{
            ID:        h.SiteID,
            Status:    string(h.Status),
            Code:      *h.StatusCode,
            CheckedAt: *h.LastCheckedAt,
        })
    }

    response := PaginatedResponse[SiteCheckHistoryItem]{
        Data:   items,
        Total:  total,
        Limit:  limit,
        Offset: offset,
    }

    writeJSON(w, http.StatusOK, response)
}
