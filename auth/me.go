package auth

import (
	"ekhoes-server/utils"
	"encoding/json"
	"net/http"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := CheckAuthorization(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Get session id

	sessionId, _ := claims["sessionId"].(string)

	// Retrieve session

	session, err := GetSession(sessionId)

	if err == SessionNotFound {
		utils.Error("Session not found")
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session.User)
}
