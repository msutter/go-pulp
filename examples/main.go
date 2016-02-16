package main

import (
	"fmt"
	"github.com/msutter/go-pulp/pulp"
	"log"
	"time"
)

func main() {
	apiUser := "admin"
	apiPasswd := "admin"
	apiEndpoint := "pulp-lab-11"

	// create the client
	client, err := pulp.NewClient(apiEndpoint, apiUser, apiPasswd, nil)

	// repository options
	ro := &pulp.GetRepositoryOptions{
		Details: true,
	}

	repo := "test-repo-1-lab"

	// get the repo
	r, _, err := client.Repositories.GetRepository(repo, ro)
	fmt.Printf("%v\n", r)
	if err != nil {
		log.Fatal(err)
	}

	// get the repo
	callReport, _, err := client.Repositories.SyncRepository(repo)
	syncTaskId := callReport.SpawnedTasks[0].TaskId
	fmt.Printf("TaskId: %v\n", syncTaskId)
	if err != nil {
		log.Fatal(err)
	}

	state := "init"
	for state != "finished" {
		task, _, terr := client.Tasks.GetTask(syncTaskId)
		fmt.Printf("%v\n", task.State)

		state = task.State
		time.Sleep(500 * time.Millisecond)
		if terr != nil {
			log.Fatal(terr)
		}
	}
}
