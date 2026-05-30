package auth

import (
	"ekhoes-server/cache"
	"ekhoes-server/config"
	"ekhoes-server/session"
	"ekhoes-server/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func createGuestSession(moduleId string, credentials Credentials, remoteAddr string) (session.Session, error) {
	utils.Debug("Creating guest session")

	user := session.User{
		Id:      utils.UUID(),
		Name:    "Guest",
		IsGuest: true,
		IsUSer:  false,
	}

	ses := session.Session{
		User:       user,
		Agent:      credentials.Agent,
		Platform:   credentials.Platform,
		Model:      credentials.Model,
		DeviceName: credentials.DeviceName,
		DeviceType: credentials.DeviceType,
		Ip:         remoteAddr,
	}

	sessionNew, err := session.Create(moduleId, ses, time.Duration(config.TTL_Session())*time.Minute)

	if err != nil {
		return ses, err
	}

	utils.Debug("Session created: %s", sessionNew.Id)

	return sessionNew, nil
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
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

	sessionId := ""
	token := ""
	refreshToken := ""
	var ses session.Session

	// Get client info

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if credentials.AppId == "" {
		http.Error(w, "Missing application id", http.StatusBadRequest)
		return
	}

	// Check token

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		token = authHeader[7:]
	}

	if token == "" {
		// Create guest session

		utils.Debug("Client doesn't have a token")

		ses, err = createGuestSession(credentials.AppId, credentials, r.RemoteAddr)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create access token

		token, err = generateAccessTokenFromSession(ses)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create refresh token

		refreshToken, err = generateRefreshToken(ses)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		// Decode token

		utils.Debug("Decoding token")

		claims, valid, err := DecodeJWT(token)

		if err != nil {
			utils.Err(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Get session id

		sessionId, _ = claims["sessionId"].(string)

		// Retrieve session

		ses, err = session.Get(sessionId)

		if errors.Is(err, session.SessionNotFound) {
			utils.Error("Session not found")
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			utils.Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if !valid {
				utils.Error("Token expired")
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			utils.Debug("Session found. Extending TTL...")

			cache.UpdateTTL(sessionId, time.Duration(config.TTL_Session())*time.Minute)

			/*
				if !valid {
					// Regenerate token

					utils.Debug("Regenerating token...")

					newClaims := CustomClaims{
						SessionId: sessionId,
						UserId:    sess.User.Id,
						Email:     credentials.Email,
						Name:      sess.User.Name,
						IsUser:    sess.User.IsUSer,
						IsGuest:   sess.User.IsGuest,
					}

					token, err = GenerateJWT(newClaims, time.Now().Add(time.Minute))

					if err != nil {
						utils.Err(err)
						http.Error(w, "Error regenerating token", http.StatusInternalServerError)
						return
					}
				}
			*/
		}
	}

	data := fmt.Sprintf(
		`{"token":"%s", "refreshToken":"%s", "name":"%s", "isGuest":%t, "isUser":%t }`,
		token,
		refreshToken,
		ses.User.Name,
		ses.User.IsGuest,
		ses.User.IsUSer)

	utils.Debug("%s", data)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
	w.WriteHeader(http.StatusOK)
}
