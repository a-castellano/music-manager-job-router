// +build integration_tests

package manager

import (
	"log"
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

func TestSendDie(t *testing.T) {

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

	var job commontypes.Job

	job.ID = "dassa111a"
	job.Status = true
	job.Finished = false
	job.Type = commontypes.Die
	job.LastOrigin = "JobManager"

	encodedJob, _ := commontypes.EncodeJob(job)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel in TestSendDie")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		testConfig.JobManager.Name, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue in TestSendDie")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	wrapperChannel := make(chan commontypes.Job)

	jobManagementError := ReadJobManagerJobs(testConfig, wrapperChannel)
	if jobManagementError != nil {
		t.Errorf("ReadJobManagerJobs should return no errors when die is processed.")
	}

	firstResultJob := <-wrapperChannel
	secondResultJob := <-wrapperChannel
	if firstResultJob.ID != job.ID {
		t.Errorf("Original and result Jobs should have same ID.")
	}
	if firstResultJob.Type != commontypes.Die {
		t.Errorf("Result Jobs type should be Die.")
	}
	expectedFirstOrigin := "first"
	if firstResultJob.RequiredOrigin != expectedFirstOrigin {
		t.Errorf("Result Job RequiredOrigin whould be '%s', not '%s'.", expectedFirstOrigin, firstResultJob.RequiredOrigin)
	}

	if secondResultJob.ID != job.ID {
		t.Errorf("Original and result Jobs should have same ID.")
	}
	if secondResultJob.Type != commontypes.Die {
		t.Errorf("Result Jobs type should be Die.")
	}
	expectedSecondOrigin := "second"
	if secondResultJob.RequiredOrigin != expectedSecondOrigin {
		t.Errorf("Second Result Job RequiredOrigin whould be '%s', not '%s'.", expectedSecondOrigin, secondResultJob.RequiredOrigin)
	}

}

func TestSendJobFromInvalidOrigin(t *testing.T) {

	var testConfig config.Config

	testConfig.Server.Host = "rabbitmq"
	testConfig.Server.Port = 5672
	testConfig.Server.User = "guest"
	testConfig.Server.Password = "guest"
	testConfig.JobManager.Name = "JobManager"

	var job commontypes.Job

	job.ID = "dassa111a"
	job.Status = true
	job.Finished = false
	job.Type = commontypes.ArtistInfoRetrieval
	job.LastOrigin = "Foo"

	encodedJob, _ := commontypes.EncodeJob(job)

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel in TestSendDie")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		testConfig.JobManager.Name, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue in TestSendDie")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	wrapperChannel := make(chan commontypes.Job)

	jobManagementError := ReadJobManagerJobs(testConfig, wrapperChannel)
	if jobManagementError != nil {
		t.Errorf("ReadJobManagerJobs should return no errors although origin is invalid.")
	}

	resultJob := <-wrapperChannel
	if resultJob.ID != job.ID {
		t.Errorf("Original and result Jobs should have same ID.")
	}
	requiredError := "LastOrigin can only be 'JobManager'"
	if resultJob.Error != requiredError {
		t.Errorf("Result Job error should be '%s', not '%s'.", requiredError, resultJob.Error)
	}
}
