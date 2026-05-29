package admin

import (
	"encoding/json"
	"log"
	"net/http"

	"ekhoes-server/auth"
	"ekhoes-server/session"

	"github.com/go-chi/chi/v5"
)

/**
 * GET /sessions
 */
func GetSessionsHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := auth.CheckAuthorization(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if auth.HasPrivilege(claims["privileges"].(string), "ek_read_session") == false {
		http.Error(w, "missing required privileges", http.StatusUnauthorized)
		return
	}

	sessions, err := session.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
	w.WriteHeader(http.StatusOK)
}

/**
 * DELETE /session/[id]
 */
func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := auth.CheckAuthorization(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if auth.HasPrivilege(claims["privileges"].(string), "ek_delete_session") == false {
		http.Error(w, "missing required privileges", http.StatusUnauthorized)
		return
	}

	sessionId := chi.URLParam(r, "id")

	err = session.Delete(sessionId)

	if err == nil {
		log.Printf("Session deleted: %s\n", sessionId)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

/**
 * DELETE /sessions
 */
func DeleteAllSessionsHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := auth.CheckAuthorization(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if auth.HasPrivilege(claims["privileges"].(string), "ek_delete_session") == false {
		http.Error(w, "missing required privileges", http.StatusUnauthorized)
		return
	}

	err = session.DeleteAll()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("All sessions deleted")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
