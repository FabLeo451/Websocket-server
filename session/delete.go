package session

import (
	"ekhoes-server/db"
	"errors"
	"fmt"
)

func Delete(sessionId string) error {
	deleted, err := db.DeleteKey(sessionId)
	if err != nil {
		return fmt.Errorf("unable to remove key: %w", err)
	}

	if !deleted {
		msg := fmt.Sprintf("session key not found for deletion: %s", sessionId)
		return errors.New(msg)
	}

	return nil
}

func DeleteAll() error {
	err := db.DeleteByPattern("*")
	if err != nil {
		return fmt.Errorf("unable to remove key: %w", err)
	}

	return nil
}
