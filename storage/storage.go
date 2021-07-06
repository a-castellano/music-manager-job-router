package storage

import (
	"encoding/json"
	"errors"
	"net/http"

	commontypes "github.com/a-castellano/music-manager-common-types/types"

	"bytes"
)

func SendInfoToStorageManager(client http.Client, storageService string, job commontypes.Job) error {

	jsonJob, _ := json.Marshal(job)
	url := "http://" + storageService
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonJob))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		// We must include a reason for that error
		return errors.New("Failed to update status.")
	}

	return nil
}
