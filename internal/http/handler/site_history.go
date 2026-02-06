package handler

import (
    "net/http"
    "strconv"
    "log/slog"

    "site-monitor/internal/domain"
)

const (
    DefaultLimit = 20
    MaxLimit     = 100
)

func (h *SiteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siteID := r.PathValue("id")      
	if siteID == "" {
		http.Error(w, "missing site id", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	limit := DefaultLimit
	offset := 0

	if l := query.Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val <= 0 {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}

		if val > MaxLimit {
			val = MaxLimit
		}

		limit = val
	}

	if o := query.Get("offset"); o != "" {
		val, err := strconv.Atoi(o)
		if err != nil || val < 0 {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}

		offset = val
	}

	history, total, err := h.checkResultRepo.GetBySiteID(ctx, siteID, limit, offset)
	if err != nil {
		h.logger.Error(
			"failed to get site history",
			slog.String("siteID", siteID),
			slog.Any("err", err),
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	items := make([]domain.CheckResult, 0, len(history))

	for _, r := range history {
		items = append(items, domain.CheckResult{
			ID: r.ID,
			HTTPStatus: r.HTTPStatus,
			IsAvailable: r.IsAvailable,
            ResponseTime: r.ResponseTime,
			CheckedAt: r.CheckedAt,
		})
	}

	response := domain.PaginatedResponse[domain.CheckResult]{
		Data:   items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	writeJSON(w, http.StatusOK, response)
}

