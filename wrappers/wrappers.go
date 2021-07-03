package wrappers

import (
	"fmt"
	"strconv"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/streadway/amqp"
)

func RouteJobs(config config.Config, wrapperChannel chan commontypes.Job) error {

	wrapperQueues := make(map[string]amqp.Queue)
	var wrapperOrder []string

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
		wrapperOrder = append(wrapperOrder, wrapper.Name)
	}

	for {
		jobToRoute := <-wrapperChannel
		encodedJob, _ := commontypes.EncodeJob(jobToRoute)
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
						return fmt.Errorf("Failed to send job to qeue %s in RouteJobs: %w", wrapperOrder[0], err)
					}

				} //else{
				//RequiredOrigin does not exists
				//}
			}
		} else {
			// Job already routed of Die sent
		}
	}

	return nil
}
