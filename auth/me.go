package auth

import (
	"ekhoes-server/session"
	"ekhoes-server/utils"
	"encoding/json"
	"errors"
	"net/http"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := CheckAuthorization(r)

	if err != nil {
		utils.Err(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get session id

	sessionId, _ := claims["sessionId"].(string)

	// Retrieve session

	ses, err := session.Get(sessionId)

	if errors.Is(err, session.SessionNotFound) {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ses.User)
}
