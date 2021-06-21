package status

import (
	"encoding/json"
	"errors"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"net/http"

	"bytes"
)

func UpdateJobStatus(client http.Client, statusService string, job commontypes.Job) error {

	jsonJob, _ := json.Marshal(job)
	url := "http://" + statusService
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
