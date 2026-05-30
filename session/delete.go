package session

import (
	"ekhoes-server/cache"
	"errors"
	"fmt"
)

func Delete(sessionId string) error {
	deleted, err := cache.DeleteKey(sessionId)
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
	err := cache.DeleteByPattern("*")
	if err != nil {
		return fmt.Errorf("unable to remove key: %w", err)
	}

	return nil
}
