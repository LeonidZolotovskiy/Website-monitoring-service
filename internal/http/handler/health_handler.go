package handler

import (
	"net/http"
	"time"
)

type HealthHandler struct {
	startTime time.Time
	version   string
}

func NewHealthHandler(startTime time.Time, version string) *HealthHandler {
	return &HealthHandler{
		startTime: startTime,
		version:   version,
	}
}

func (h *HealthHandler) checkDependencies() []DependencyStatus {
	// Сейчас зависимостей нет — возвращаем пустой список
	// В будущем сюда добавится БД, Redis, внешние API и т.д.
	return []DependencyStatus{}
}
// Check godoc
// @Summary      Health check
// @Description  Returns service health status, uptime, version and dependencies state
// @Tags         health
// @Produce      json
// @Success      200 {object} HealthResponse "Service is healthy"
// @Failure      503 {object} HealthResponse "Service is unhealthy"
// @Router       /health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	uptime := now.Sub(h.startTime)

	dependencies := h.checkDependencies()

	status := HealthHealthy
	httpStatus := http.StatusOK

	for _, d := range dependencies {
		if d.Status != "healthy" {
			status = HealthUnhealthy
			httpStatus = http.StatusServiceUnavailable
			break
		}
	}

	resp := HealthResponse{
		Status:       status,
		Version:      h.version,
		Uptime:       uptime.Round(time.Second).String(),
		ServerTime:   now,
		Dependencies: dependencies,
	}

	writeJSON(w, httpStatus, resp)
}
