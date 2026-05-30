package auth

import (
	"ekhoes-server/cache"
	"ekhoes-server/session"
	"encoding/json"
	"errors"
	"net/http"
)

type RefreshData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var data RefreshData

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check refresh token exists

	sessionId, err := cache.Get(data.RefreshToken)

	if err == cache.KeyNotFound {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get session

	ses, err := session.Get(sessionId)

	if errors.Is(err, session.SessionNotFound) {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate access token

	data.AccessToken, err = generateAccessTokenFromSession(ses)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rotate refresh token

	data.RefreshToken, err = rotateRefreshToken(data.RefreshToken, ses)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
