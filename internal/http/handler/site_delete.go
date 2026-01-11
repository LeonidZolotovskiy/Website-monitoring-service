package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"site-monitor/internal/repository"
)

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