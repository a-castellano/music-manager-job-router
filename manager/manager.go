package manager

import (
	"errors"
	"fmt"
	"strconv"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/streadway/amqp"
)

func ReadJobManagerJobs(config config.Config, wrapperChannel chan commontypes.Job) error {

	connection_string := "amqp://" + config.Server.User + ":" + config.Server.Password + "@" + config.Server.Host + ":" + strconv.Itoa(config.Server.Port) + "/"
	conn, err := amqp.Dial(connection_string)

	if err != nil {
		return fmt.Errorf("Failed to stablish connection with RabbitMQ: %w", err)
	}
	defer conn.Close()

	jobmanager_ch, err := conn.Channel()
	defer jobmanager_ch.Close()

	if err != nil {
		return fmt.Errorf("Failed to open jobmanager RabbitMQ channel: %w", err)
	}

	jobmanager_q, err := jobmanager_ch.QueueDeclare(
		config.JobManager.Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		return fmt.Errorf("Failed to declare jobmanager queue: %w", err)
	}

	err = jobmanager_ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return fmt.Errorf("Failed to set jobmanager QoS: %w", err)
	}

	jobsToProcess, err := jobmanager_ch.Consume(
		jobmanager_q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %w", err)
	}

	processJobs := make(chan bool)

	go func() {
		for job := range jobsToProcess {

			die := false
			jobToProcess, decodeJobErr := commontypes.DecodeJob(job.Body)

			if decodeJobErr != nil {
				err = errors.New("Empty job data received.")
			} else {
				if jobToProcess.Type == commontypes.Die {
					die = true
				} else {
					// This function  reads meesages from jobManager
					if jobToProcess.LastOrigin != "JobManager" {
						jobToProcess.Error = "LastOrigin can only be 'JobManager'"
						jobToProcess.Status = false
						job.Ack(false)
						processJobs <- true
						wrapperChannel <- jobToProcess
					} else {
						jobToProcess.LastOrigin = "JobRouter"
						processJobs <- true
						wrapperChannel <- jobToProcess
					}
				}
			}
			if die {
				job.Ack(false)
				processJobs <- false
				for _, wrapper := range config.Wrappers {
					jobToWrapper := jobToProcess
					jobToWrapper.RequiredOrigin = wrapper.Name
					wrapperChannel <- jobToWrapper
				}
				return
			}
		}
		return
	}()

	<-processJobs

	return nil
}
