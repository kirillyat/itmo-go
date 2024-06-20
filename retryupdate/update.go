//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(client kvapi.Client, key string, updateFn func(currentValue *string) (newValue string, err error)) error {
	var currentValue *string
	var currentVersion uuid.UUID
	var authError *kvapi.AuthError
	for isVersionAcquired := false; !isVersionAcquired; {
		getResponse, err := client.Get(&kvapi.GetRequest{Key: key})

		switch true {
		case errors.Is(err, kvapi.ErrKeyNotFound):
			isVersionAcquired = true
		case err == nil:
			isVersionAcquired = true
			currentValue = &getResponse.Value
			currentVersion = getResponse.Version
		case errors.As(err, &authError):
			return err
		}
	}

	updatedValue, updateErr := updateFn(currentValue)
	if updateErr != nil {
		return updateErr
	}
	var conflictError *kvapi.ConflictError
	newVersion := uuid.Must(uuid.NewV4())
	for isWriteSuccessful := false; !isWriteSuccessful; {
		_, setErr := client.Set(&kvapi.SetRequest{Key: key, Value: updatedValue, OldVersion: currentVersion, NewVersion: newVersion})

		switch true {
		case errors.Is(setErr, kvapi.ErrKeyNotFound):
			currentVersion = uuid.UUID{}
			updatedValue, updateErr = updateFn(nil)
			if updateErr != nil {
				return updateErr
			}
		case setErr == nil || errors.As(setErr, &authError):
			return setErr
		case errors.As(setErr, &conflictError):
			if conflictError.ExpectedVersion == newVersion {
				return nil
			}
			return UpdateValue(client, key, updateFn)
		}
	}

	return nil
}
