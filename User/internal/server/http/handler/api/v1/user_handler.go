package v1

import (
	"User/internal/domain/user"
	"User/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	svc       service.UserService
	url       service.UrlService
	jwtSecret []byte
}

func NewUserHandler(svc service.UserService, url service.UrlService, jwtSecret []byte) *UserHandler {
	return &UserHandler{svc: svc, url: url, jwtSecret: jwtSecret}
}

// respondWithJSON helper
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError helper
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RegisterHandler handles POST /api/v1/user/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	uid, err := h.svc.Register(r.Context(), &u)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int64{"uid": uid})
}

// LoginHandler handles POST /api/v1/user/login (The actual login)
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.svc.Login(r.Context(), creds.Username, creds.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// LoginCheckHandler handles GET /api/v1/user/login (JWT Check)
func (h *UserHandler) LoginCheck(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing auth token")
		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	uid := int64(claims["uid"].(float64))
	u, err := h.svc.GetUserByID(r.Context(), uid)
	if err != nil || u == nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

// GetUserUrlsHandler handles GET /api/v1/user/{userid}?count=xx
func (h *UserHandler) GetUserUrls(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userid")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	countStr := r.URL.Query().Get("count")
	count, _ := strconv.Atoi(countStr)
	if count <= 0 {
		count = 10
	}

	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	urls, err := h.url.GetUserUrls(r.Context(), userID, count, page)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, urls)
}
