// +build integration_tests

package wrappers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestReceiveDie(t *testing.T) {

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)

	var job commontypes.Job

	job.ID = "dassa111a"
	job.Status = true
	job.Finished = false
	job.Type = commontypes.Die
	job.LastOrigin = "JobRouter"
	job.RequiredOrigin = "JobRouter"

	wrapperChannel := make(chan commontypes.Job)

	go func() { wrapperChannel <- job }()

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	err := RouteJobs(testConfig, wrapperChannel, client)

	if err != nil {
		t.Errorf("TestReceiveDie should end without errors.")
	}

}

func TestReceiveFinishedJobAndDie(t *testing.T) {

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)

	var dieJob commontypes.Job
	var finishedJob commontypes.Job

	dieJob.ID = "dassa111a"
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die
	dieJob.LastOrigin = "JobRouter"
	dieJob.RequiredOrigin = "JobRouter"

	finishedJob.ID = "dassa111a"
	finishedJob.Status = true
	finishedJob.Finished = true
	finishedJob.Type = commontypes.ArtistInfoRetrieval
	finishedJob.LastOrigin = "first"

	wrapperChannel := make(chan commontypes.Job)

	go func() {
		wrapperChannel <- finishedJob
		wrapperChannel <- dieJob
	}()

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	err := RouteJobs(testConfig, wrapperChannel, client)

	if err != nil {
		t.Errorf("TestReceiveFinishedJobAndDie should end without errors.")
	}

}

func TestReceiveFailedJobNoMoreWrappersJobAndDie(t *testing.T) {

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)

	var dieJob commontypes.Job
	var unfinishedJob commontypes.Job

	dieJob.ID = "dassa111a"
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die
	dieJob.LastOrigin = "JobRouter"
	dieJob.RequiredOrigin = "JobRouter"

	unfinishedJob.ID = "dassa111a"
	unfinishedJob.Status = false
	unfinishedJob.Finished = false
	unfinishedJob.Type = commontypes.ArtistInfoRetrieval
	unfinishedJob.LastOrigin = "first"

	wrapperChannel := make(chan commontypes.Job)

	go func() {
		wrapperChannel <- unfinishedJob
		wrapperChannel <- dieJob
	}()

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	err := RouteJobs(testConfig, wrapperChannel, client)

	if err != nil {
		t.Errorf("TestReceiveFailedJobNoMoreWrappersJobAndDie should end without errors.")
	}

}

func TestReceiveFailedJobOneMoreWrapperJobAndDie(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}
	secondwrapper := config.Queue{Name: "second"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)
	testConfig.Wrappers = append(testConfig.Wrappers, secondwrapper)

	var dieJob commontypes.Job
	var unfinishedJob commontypes.Job

	dieJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die
	dieJob.LastOrigin = "JobRouter"
	dieJob.RequiredOrigin = "JobRouter"

	unfinishedJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	unfinishedJob.Status = false
	unfinishedJob.Finished = false
	unfinishedJob.Type = commontypes.ArtistInfoRetrieval
	unfinishedJob.LastOrigin = "first"

	wrapperChannel := make(chan commontypes.Job)

	go func() {
		wrapperChannel <- unfinishedJob
		wrapperChannel <- dieJob
	}()

	err := RouteJobs(testConfig, wrapperChannel, client)

	connectionString := "amqp://" + testConfig.Server.User + ":" + testConfig.Server.Password + "@" + testConfig.Server.Host + ":" + strconv.Itoa(testConfig.Server.Port) + "/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		failOnError(err, "Failed to stablish connection with RabbitMQ")
	}
	defer conn.Close()

	secondWrapperCh, err := conn.Channel()
	defer secondWrapperCh.Close()

	if err != nil {
		failOnError(err, "Failed to open first Wrapper RabbitMQ channel")
	}

	secondWrapperQ, err := secondWrapperCh.QueueDeclare(
		testConfig.Wrappers[1].Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare secondWrapper queue.")
	}

	err = secondWrapperCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		failOnError(err, "Failed to set firstWrapper QoS.")
	}

	jobsToProcess, err := secondWrapperCh.Consume(
		secondWrapperQ.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)
	var receivedData []byte

	go func() {
		for d := range jobsToProcess {

			receivedData = d.Body
			d.Ack(true)
			forever <- false
		}
	}()

	<-forever
	decodedJob, decodedJobErr := commontypes.DecodeJob(receivedData)

	if decodedJobErr != nil {
		t.Errorf("TestReceive shouldn't fail decoding job.")
	}
	if decodedJob.ID != unfinishedJob.ID {
		t.Errorf("job and decodedJob should have the same ID.")
	}
}

func TestReceiveJobFirstWrapperJobAndDie(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}
	secondwrapper := config.Queue{Name: "second"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)
	testConfig.Wrappers = append(testConfig.Wrappers, secondwrapper)

	var dieJob commontypes.Job
	var unfinishedJob commontypes.Job

	dieJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die
	dieJob.LastOrigin = "JobRouter"
	dieJob.RequiredOrigin = "JobRouter"

	unfinishedJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	unfinishedJob.Status = true
	unfinishedJob.Finished = false
	unfinishedJob.Type = commontypes.ArtistInfoRetrieval
	unfinishedJob.LastOrigin = "JobManager"

	wrapperChannel := make(chan commontypes.Job)

	go func() {
		wrapperChannel <- unfinishedJob
		wrapperChannel <- dieJob
	}()

	err := RouteJobs(testConfig, wrapperChannel, client)

	connectionString := "amqp://" + testConfig.Server.User + ":" + testConfig.Server.Password + "@" + testConfig.Server.Host + ":" + strconv.Itoa(testConfig.Server.Port) + "/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		failOnError(err, "Failed to stablish connection with RabbitMQ")
	}
	defer conn.Close()

	secondWrapperCh, err := conn.Channel()
	defer secondWrapperCh.Close()

	if err != nil {
		failOnError(err, "Failed to open first Wrapper RabbitMQ channel")
	}

	secondWrapperQ, err := secondWrapperCh.QueueDeclare(
		testConfig.Wrappers[0].Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare secondWrapper queue.")
	}

	err = secondWrapperCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		failOnError(err, "Failed to set firstWrapper QoS.")
	}

	jobsToProcess, err := secondWrapperCh.Consume(
		secondWrapperQ.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)
	var receivedData []byte

	go func() {
		for d := range jobsToProcess {

			receivedData = d.Body
			d.Ack(true)
			forever <- false
		}
	}()

	<-forever
	decodedJob, decodedJobErr := commontypes.DecodeJob(receivedData)

	if decodedJobErr != nil {
		t.Errorf("TestReceive shouldn't fail decoding job.")
	}
	if decodedJob.ID != unfinishedJob.ID {
		t.Errorf("job and decodedJob should have the same ID.")
	}
}

func TestReceiveJobRequiredOriginFirstWrapperJobAndDie(t *testing.T) {

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(`
	not html code
		`))}}}

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	firstwrapper := config.Queue{Name: "first"}
	secondwrapper := config.Queue{Name: "second"}

	testConfig.Wrappers = append(testConfig.Wrappers, firstwrapper)
	testConfig.Wrappers = append(testConfig.Wrappers, secondwrapper)

	var dieJob commontypes.Job
	var unfinishedJob commontypes.Job

	dieJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die
	dieJob.LastOrigin = "JobRouter"
	dieJob.RequiredOrigin = "JobRouter"

	unfinishedJob.ID = "TestReceiveFailedJobOneMoreWrapperJobAndDie"
	unfinishedJob.Status = true
	unfinishedJob.Finished = false
	unfinishedJob.Type = commontypes.ArtistInfoRetrieval
	unfinishedJob.LastOrigin = "JobManager"
	unfinishedJob.RequiredOrigin = "second"

	wrapperChannel := make(chan commontypes.Job)

	go func() {
		wrapperChannel <- unfinishedJob
		wrapperChannel <- dieJob
	}()

	err := RouteJobs(testConfig, wrapperChannel, client)

	connectionString := "amqp://" + testConfig.Server.User + ":" + testConfig.Server.Password + "@" + testConfig.Server.Host + ":" + strconv.Itoa(testConfig.Server.Port) + "/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		failOnError(err, "Failed to stablish connection with RabbitMQ")
	}
	defer conn.Close()

	secondWrapperCh, err := conn.Channel()
	defer secondWrapperCh.Close()

	if err != nil {
		failOnError(err, "Failed to open first Wrapper RabbitMQ channel")
	}

	secondWrapperQ, err := secondWrapperCh.QueueDeclare(
		testConfig.Wrappers[1].Name, //second
		true,                        // Durable
		false,                       // DeleteWhenUnused
		false,                       // Exclusive
		false,                       // NoWait
		nil,                         // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare secondWrapper queue.")
	}

	err = secondWrapperCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		failOnError(err, "Failed to set firstWrapper QoS.")
	}

	jobsToProcess, err := secondWrapperCh.Consume(
		secondWrapperQ.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		failOnError(err, "Failed to register a consumer")
	}

	forever := make(chan bool)
	var receivedData []byte

	go func() {
		for d := range jobsToProcess {

			receivedData = d.Body
			d.Ack(true)
			forever <- false
		}
	}()

	<-forever
	decodedJob, decodedJobErr := commontypes.DecodeJob(receivedData)

	if decodedJobErr != nil {
		t.Errorf("TestReceive shouldn't fail decoding job.")
	}
	if decodedJob.ID != unfinishedJob.ID {
		t.Errorf("job and decodedJob should have the same ID.")
	}
}
