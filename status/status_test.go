// +build integration_tests unit_tests

package status

import (
	"bytes"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	uuid "github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestUpdateJobStatusFailedStatusCode(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 500, Body: ioutil.NopCloser(bytes.NewBufferString(`
not html code
	`))}}}

	var statusServiceName string = "Test"

	var newJob commontypes.Job
	newJob.ID = uuid.New()
	newJob.Type = commontypes.RecordInfoRetrieval

	err := UpdateJobStatus(client, statusServiceName, newJob)

	if err == nil {
		t.Errorf("UpdateJobStatusFailedStatusCode should fail.")
	}
}

func TestUpdateJobStatusSuccessStatusCode(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
not html code
	`))}}}

	var statusServiceName string = "Test"

	var newJob commontypes.Job
	newJob.ID = uuid.New()
	newJob.Type = commontypes.RecordInfoRetrieval

	err := UpdateJobStatus(client, statusServiceName, newJob)

	if err != nil {
		t.Errorf("UpdateJobStatusFailedStatusCode shouldn't fail.")
	}
}
