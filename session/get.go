package session

import (
	"ekhoes-server/db"
	"encoding/json"
	"log"
)

func Get(id string) (Session, error) {
	val, err := db.Get(id)

	var sess Session

	if err == db.KeyNotFound {
		return sess, SessionNotFound
	}

	err = json.Unmarshal([]byte(val), &sess)

	sess.TTL = db.GetTTL(id)

	return sess, err
}

func GetAll() ([]Session, error) {
	var sessions []Session

	keys, err := db.GetKeysByPattern("ses:*")
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		val, err := db.Get(key)
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
		sess.TTL = db.GetTTL(key)

		sessions = append(sessions, sess)
	}

	return sessions, nil
}
