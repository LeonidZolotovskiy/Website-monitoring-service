package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"site-monitor/internal/repository"
)
// Delete godoc
// @Summary      Delete site
// @Description  Deletes a site by ID
// @Tags         sites
// @Param        id path string true "Site ID"
// @Success      204 "Site deleted"
// @Failure      404 {string} string "Site not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /sites/{id} [delete]
func (h *SiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)      
    id := vars["id"]        

    if err := h.siteRepo.DeleteByID(id); err != nil {
        if err == repository.ErrSiteNotFound {
            http.Error(w, "site not found", http.StatusNotFound)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}