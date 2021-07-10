package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	commontypes "github.com/a-castellano/music-manager-common-types/types"
	"github.com/a-castellano/music-manager-job-router/config"
	"github.com/a-castellano/music-manager-job-router/manager"
	"github.com/a-castellano/music-manager-job-router/wrappers"
)

func main() {

	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	log.Println("Reading config.")

	jobRouterConfig, err := config.ReadConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		log.Println("Config readed successfully.")

		wrapperChannel := make(chan commontypes.Job)
		go manager.ReadJobManagerJobs(jobRouterConfig, wrapperChannel)
		jobRouterError := wrappers.RouteJobs(jobRouterConfig, wrapperChannel, client)

		if jobRouterError != nil {
			fmt.Println(jobRouterError)
		}
	}

}
