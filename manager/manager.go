package manager

import (
	"fmt"
	"strconv"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/streadway/amqp"
)

func ReadJobManagerJobs(config config.Config, wrapper_channel chan commontypes.Job) error {

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
		return fmt.Errorf("Failed to declare incoming queue: %w", err)
	}

	return nil
}
