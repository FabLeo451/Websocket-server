package auth

import (
	"ekhoes-server/db"
	"encoding/json"
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

	sessionId, err := db.Get(data.RefreshToken)

	if err == db.KeyNotFound {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get session

	session, err := GetSession(sessionId)

	if err == SessionNotFound {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate access token

	data.AccessToken, err = generateAccessTokenFromSession(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rotate refresh token

	data.RefreshToken, err = rotateRefreshToken(data.RefreshToken, session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
