package main

import (
	"fmt"
	"github.com/msutter/go-pulp/pulp"
	"log"
)

func main() {
	apiUser := "admin"
	apiPasswd := "admin"
	apiEndpoint := "pulp-lab-11"
	client := pulp.NewClient(apiEndpoint, apiUser, apiPasswd, nil)

	// List all repos
	repos, _, err := client.Repositories.ListRepositories()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", repos)

	// Create a new repository file
	ro := &pulp.CreateRepositoryOptions{
		Details: true,
	}

	for _, repo := range repos {
		r, _, err := client.Repositories.GetRepository(repo.Id, ro)
		fmt.Printf("%v\n", r)
		if err != nil {
			log.Fatal(err)
		}
	}

	taskId := "bac41a2b-0830-4038-8bb9-2d917624b888"
	task, _, terr := client.Tasks.GetTask(taskId)
	fmt.Printf("%v\n", task)
	if terr != nil {
		log.Fatal(terr)
	}

	_ = "breakpoint"

}
