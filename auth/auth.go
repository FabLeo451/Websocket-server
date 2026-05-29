package auth

import (
	"ekhoes-server/config"
	"ekhoes-server/db"
	"ekhoes-server/session"
	"ekhoes-server/utils"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credentials struct {
	AppId      string `json:"appId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Agent      string `json:"agent"`
	Platform   string `json:"platform"`
	Model      string `json:"model"`
	DeviceName string `json:"deviceName"`
	DeviceType string `json:"deviceType"`
}

func CheckAuthorization(r *http.Request) (jwt.MapClaims, error) {

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	// expected format:
	// Authorization: Bearer <token>

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid authorization header format")
	}

	scheme := parts[0]
	token := strings.TrimSpace(parts[1])

	if !strings.EqualFold(scheme, "Bearer") {
		return nil, errors.New("invalid authorization scheme")
	}

	if token == "" {
		return nil, errors.New("missing bearer token")
	}

	claims, valid, err := DecodeJWT(token)

	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors.New("token expired")
	}

	/*
		userId, ok := claims["userId"].(string)

		if !ok || userId == "" {
			return nil, errors.New("missing user id in token")
		}
	*/

	return claims, nil
}

func contains(csv string, target string) bool {
	items := strings.Split(csv, ",")
	for _, item := range items {
		if strings.TrimSpace(item) == target {
			return true
		}
	}
	return false
}

func HasPrivilege(privileges string, target string) bool {
	return contains(privileges, target) || contains(privileges, "ek_admin")
}

func generateAccessTokenFromSession(session session.Session) (string, error) {
	claims := CustomClaims{
		SessionId: session.Id,
		IsUser:    session.User.IsUSer,
		IsGuest:   session.User.IsUSer,
	}

	token, err := GenerateJWT(claims, time.Now().Add(time.Duration(config.TTL_Token())*time.Minute))

	if err != nil {
		return "", err
	}

	return token, nil
}

func generateRefreshToken(session session.Session) (string, error) {

	refreshToken := "rt_" + utils.RandomId()

	err := db.SetWithTTL(refreshToken, []byte(session.Id), time.Duration(config.TTL_RefreshToken())*time.Minute)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func rotateRefreshToken(oldToken string, session session.Session) (string, error) {
	_, err := db.DeleteKey(oldToken)

	if err != nil {
		return "", err
	}

	refreshToken, err := generateRefreshToken(session)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
