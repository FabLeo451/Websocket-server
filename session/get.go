package session

import (
	"ekhoes-server/cache"
	"encoding/json"
	"log"
)

func Get(id string) (Session, error) {
	val, err := cache.Get(id)

	var sess Session

	if err == cache.KeyNotFound {
		return sess, SessionNotFound
	}

	err = json.Unmarshal([]byte(val), &sess)

	sess.TTL = cache.GetTTL(id)

	return sess, err
}

func GetAll() ([]Session, error) {
	var sessions []Session

	keys, err := cache.GetKeysByPattern("ses:*")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		val, err := cache.Get(key)
		if err != nil {
			log.Printf("Errore nel leggere chiave %s: %v", key, err)
			continue
		}

		var sess Session
		if err := json.Unmarshal([]byte(val), &sess); err != nil {
			log.Printf("Errore nel parsing JSON della chiave %s: %v", key, err)
			continue
		}

		sess.Id = key
		sess.TTL = cache.GetTTL(key)

		sessions = append(sessions, sess)
	}

	return sessions, nil
}
