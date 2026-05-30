package session

import (
	"ekhoes-server/cache"
	"encoding/json"
	"errors"
	"time"

	"log"
)

type User struct {
	Id         string `json:"id" bson:"Id"`
	Name       string `json:"name" bson:"Name"`
	Email      string `json:"email" bson:"Email"`
	Roles      string `json:"roles" bson:"Roles"`
	Privileges string `json:"privileges" bson:"Privileges"`
	IsGuest    bool   `json:"isGuest" bson:"IsGuest"`
	IsUSer     bool   `json:"isUSer" bson:"IsUSer"`
}

type Session struct {
	Id         string        `json:"id"`
	User       User          `json:"user"`
	Agent      string        `json:"agent"`
	Platform   string        `json:"platform"`
	Model      string        `json:"model"`
	DeviceName string        `json:"deviceName"`
	DeviceType string        `json:"deviceType"`
	Ip         string        `json:"ip"`
	Status     string        `json:"status"`
	Created    time.Time     `json:"created"`
	Updated    time.Time     `json:"updated"`
	TTL        time.Duration `json:"ttl"`
}

var SessionNotFound = errors.New("session not found")

func SetActive(sessionId string, active bool) (Session, bool) {
	var session Session

	sessionStr, err := cache.Get(sessionId)
	if err != nil {
		return session, false
	}

	err = json.Unmarshal([]byte(sessionStr), &session)
	if err != nil {
		panic(err)
	}

	if active {
		session.Status = "online"
	} else {
		session.Status = "idle"
	}

	session.Updated = time.Now().UTC()

	modifiedJSON, err := json.Marshal(session)
	if err != nil {
		panic(err)
	}

	err = cache.Update(sessionId, modifiedJSON)

	if err != nil {
		log.Fatalf("Error updating: %v", err)
	}

	return session, true
}
