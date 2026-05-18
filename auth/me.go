package auth

import (
	"ekhoes-server/utils"
	"encoding/json"
	"net/http"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
/*
		dump, err := httputil.DumpRequest(r, true) // true = include il body
		if err != nil {
			fmt.Println("Errore DumpRequest:", err)
			return
		}

		fmt.Println("===== HTTP REQUEST DUMP =====")
		fmt.Println(string(dump))
		fmt.Println("===== END REQUEST =====")
*/

	claims, err := CheckAuthorization(r)

	if err != nil {
		utils.Err(err)
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
