package wrappers

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"strconv"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/a-castellano/music-manager-job-router/status"
	"github.com/a-castellano/music-manager-job-router/storage"
	"github.com/streadway/amqp"
)

func RouteJobs(config config.Config, wrapperChannel chan commontypes.Job, client http.Client) error {

	syslogger, err := syslog.New(syslog.LOG_INFO, "RouteJobs")
	log.SetOutput(syslogger)
	log.Println("Start")

	wrapperQueues := make(map[string]amqp.Queue)
	wrapperQueuesPosition := make(map[string]int)
	var wrapperOrder []string
	var wrapperCounter int = 0

	connection_string := "amqp://" + config.Server.User + ":" + config.Server.Password + "@" + config.Server.Host + ":" + strconv.Itoa(config.Server.Port) + "/"
	conn, err := amqp.Dial(connection_string)

	if err != nil {
		return fmt.Errorf("Failed to stablish connection with RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel in RouteJobs: %w", err)
	}
	defer ch.Close()

	for _, wrapper := range config.Wrappers {
		wrapperQueue, err := ch.QueueDeclare(
			wrapper.Name, // name
			true,         // durable
			false,        // delete when unused
			false,        // exclusive
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("Failed to declare queue %s in RouteJobs: %w", wrapper.Name, err)
		}
		wrapperQueues[wrapper.Name] = wrapperQueue
		wrapperQueuesPosition[wrapper.Name] = wrapperCounter
		wrapperOrder = append(wrapperOrder, wrapper.Name)
		wrapperCounter++
	}

	log.Println("Before for")
	for {
		jobToRoute := <-wrapperChannel
		encodedJob, _ := commontypes.EncodeJob(jobToRoute)
		log.Println("jobToRoute.LastOrigin is", jobToRoute.LastOrigin)
		if jobToRoute.LastOrigin == "JobManager" {
			if jobToRoute.RequiredOrigin == "" {
				// Send to first wrapper
				err = ch.Publish(
					"",              // exchange
					wrapperOrder[0], // routing key
					false,           // mandatory
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         encodedJob,
					})
				if err != nil {
					return fmt.Errorf("Failed to send job to qeue %s in RouteJobs: %w", wrapperOrder[0], err)
				}
			} else {
				// check if required origin exists
				if _, ok := wrapperQueues[jobToRoute.RequiredOrigin]; ok {
					err = ch.Publish(
						"",                        // exchange
						jobToRoute.RequiredOrigin, // routing key
						false,                     // mandatory
						false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         encodedJob,
						})
					if err != nil {
						return fmt.Errorf("Failed to send job to qeue %s in RouteJobs: %w", jobToRoute.RequiredOrigin, err)
					}

				} //else{
				//RequiredOrigin does not exists
				//}
			}
		} else {
			// Job has already been proccesed by another of Die signal has been sent
			log.Println("jobToRoute.Status is", jobToRoute.Status)
			if jobToRoute.Status == false {
				//Job failed - check if there are wrappers left to process this job
				if jobToRoute.RequiredOrigin == "" && wrapperQueuesPosition[jobToRoute.LastOrigin] < len(wrapperOrder) {
					// Send job to next wrapper
					nextWrapper := wrapperOrder[wrapperQueuesPosition[jobToRoute.LastOrigin]+1]
					err = ch.Publish(
						"",          // exchange
						nextWrapper, // routing key
						false,       // mandatory
						false,
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         encodedJob,
						})
					if err != nil {
						return fmt.Errorf("Failed to send job to qeue %s in RouteJobs: %w", nextWrapper, err)
					}

				} else {
					// No more wrappers left, job is marked as failed
					jobToRoute.Finished = true
					err = status.UpdateJobStatus(client, config.Status, jobToRoute)
					if err != nil {
						return fmt.Errorf("Failed to send job to status Manager in RouteJobs: %w", err)
					}
				}
			} else {
				// jobFinished or is a Die function
				log.Println("jobToRoute.RequiredOrigin is", jobToRoute.RequiredOrigin)
				if jobToRoute.RequiredOrigin == "JobRouter" {
					log.Println("jobToRoute.Type is", jobToRoute.Type)

					if jobToRoute.Type == commontypes.Die {
						break
					} else {
						return fmt.Errorf("Only JobType allowed when RequiredOrigin is JobRouter is Die.")
					}
				}
				jobToRoute.Finished = true
				err = status.UpdateJobStatus(client, config.Status, jobToRoute)
				if err != nil {
					return fmt.Errorf("Failed to send job to status Manager in RouteJobs: %w", err)
				}
				err = storage.SendInfoToStorageManager(client, config.Storage, jobToRoute)
				if err != nil {
					return fmt.Errorf("Failed to send job to status Manager in RouteJobs: %w", err)
				}

			}
		}
	}

	return nil
}
