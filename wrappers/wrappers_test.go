// +build integration_tests

package wrappers

import (
	"log"
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

func TestReceive(t *testing.T) {

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
	job.Type = commontypes.ArtistInfoRetrieval
	job.LastOrigin = "JobRouter"

	wrapperChannel := make(chan commontypes.Job)

	wrapperChannel <- job

	connectionString := "amqp://" + testConfig.Server.User + ":" + testConfig.Server.Password + "@" + testConfig.Server.Host + ":" + strconv.Itoa(testConfig.Server.Port) + "/"
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		failOnError(err, "Failed to stablish connection with RabbitMQ")
	}
	defer conn.Close()

	firstWrapperCh, err := conn.Channel()
	defer firstWrapperCh.Close()

	if err != nil {
		failOnError(err, "Failed to open first Wrapper RabbitMQ channel")
	}

	firstWrapperQ, err := firstWrapperCh.QueueDeclare(
		testConfig.Wrappers[0].Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare firstWrapper queue.")
	}

	err = firstWrapperCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		failOnError(err, "Failed to set firstWrapper QoS.")
	}

	jobsToProcess, err := firstWrapperCh.Consume(
		firstWrapperQ.Name,
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
	if decodedJob.ID != job.ID {
		t.Errorf("job and decodedJob should have the same ID.")
	}
}
