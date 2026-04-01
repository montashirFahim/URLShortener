package v1

import (
	"Server/internal/entities/auth"
	"Server/internal/entities/user"
	"Server/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	svc service.UserService
}

func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newUser, err := h.svc.Register(r.Context(), &u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":         newUser.ID,
		"username":   newUser.UserName,
		"email":      newUser.Email,
		"created_at": newUser.CreatedAt,
	})
}

func (h *AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	var req auth.TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var resp *auth.TokenResponse
	var err error

	switch req.GrantType {
	case "password":
		resp, err = h.svc.Login(r.Context(), &req)
	case "refresh_token":
		resp, err = h.svc.RefreshToken(r.Context(), &req)
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid grant type")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	var req auth.RevokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.svc.RevokeToken(r.Context(), req.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "token revoked successfully"})
}
