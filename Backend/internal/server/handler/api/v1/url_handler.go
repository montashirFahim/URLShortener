package v1

import (
	"Server/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL        string `json:"url"`
		CustomCode string `json:"custom_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	uid := r.Context().Value("uid").(string)
	su, lu, err := h.svc.ShortenURL(r.Context(), uid, req.URL, req.CustomCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":        su.ID,
		"code":      su.Code,
		"long_url":  lu.Url,
		"created_at": su.CreatedAt,
	})
}

func (h *URLHandler) GenerateGuest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	uid := r.Context().Value("uid").(string)
	su, _, err := h.svc.ShortenURL(r.Context(), uid, req.URL, "")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"code":       su.Code,
		"expires_at": su.ExpiresAt,
	})
}

func (h *URLHandler) GetGuestInfo(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "short_url")
	longURL, err := h.svc.GetLongURL(r.Context(), code)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "URL not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"long_url": longURL,
		"code":     code,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "short_url")
	longURL, err := h.svc.GetLongURL(r.Context(), code)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "URL not found")
		return
	}

	// Non-blocking analytics
	go h.svc.RecordClick(r.Context(), code, r.RemoteAddr, r.UserAgent())

	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 10
	}

	urls, total, err := h.svc.GetUserURLs(r.Context(), uid, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": urls,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(chi.URLParam(r, "url_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	detail, err := h.svc.GetURLDetail(r.Context(), uid, urlID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(chi.URLParam(r, "url_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	err = h.svc.DeleteURL(r.Context(), uid, urlID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(chi.URLParam(r, "url_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL ID")
		return
	}

	urlStats, err := h.svc.GetAnalytics(r.Context(), uid, urlID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"clicks":      urlStats.Clicks,
		"last_access": urlStats.LastAccess,
	})
}
