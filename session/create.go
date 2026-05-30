package session

import (
	"ekhoes-server/cache"
	"ekhoes-server/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func Create(appId string, session Session, ttl time.Duration) (Session, error) {
	session.Created = time.Now().UTC()
	session.Updated = time.Now().UTC()
	session.Id = fmt.Sprintf("ses:%s:%s", appId, utils.ULID())

	if session.Status == "" {
		session.Status = "idle"
	}

	data, err := json.Marshal(session)
	if err != nil {
		return session, err
	}

	err = cache.SetWithTTL(session.Id, data, ttl)

	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	return session, nil
}
