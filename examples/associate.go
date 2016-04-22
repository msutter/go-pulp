// simple example of getting a repo metadatas and running a sync on it
package main

import (
	"encoding/json"
	"fmt"
	"github.com/msutter/go-pulp/pulp"
	"log"
	"time"
)

func main() {
	apiUser := "admin"
	apiPasswd := "admin"
	apiEndpoint := "pulp-lab-11.test"

	DisableSsl := false
	SkipSslVerify := true

	// create the client
	client, err := pulp.NewClient(apiEndpoint, apiUser, apiPasswd, DisableSsl, SkipSslVerify, nil)

	target_repo := "test-repo-1-lab"
	source_repo := "sccloud-mgmt-infra-el6-lab"

	fileNames := []string{
		"nodetree",
	}

	// get units
	// Get repo units
	us, _, uerr := client.Repositories.ListRepositoryUnits(source_repo)

	if uerr != nil {
		fmt.Println(uerr.Error())
		log.Fatal(uerr)
	}

	fmt.Printf("Units in %v\n", source_repo)
	for _, u := range us {
		fmt.Printf("%v\n", u.Metadata.FileName)
	}

	orFilter := pulp.NewFilter("$or")
	for _, fileName := range fileNames {
		// only match exact filename
		regex := fmt.Sprintf("^%s", fileName)
		orFilter.AddExpression("filename", "$regex", regex)
	}

	andFilter := pulp.NewFilter("$and")
	andFilter.AddExpression("filename", "$regex", "node")

	criteria := pulp.NewUnitAssociationCriteria()
	criteria.AddFilter(orFilter)
	criteria.AddFilter(andFilter)
	criteria.AddField("name")
	criteria.AddField("version")
	criteria.AddField("epoch")
	criteria.AddField("release")
	criteria.AddField("arch")
	criteria.AddField("checksumtype")
	criteria.AddField("checksum")

	// check request body
	jsonCriteria, err := json.Marshal(criteria)
	fmt.Printf("Criteria: %s\n", jsonCriteria)

	associateCallReport, _, err := client.Repositories.CopyRepositoryUnits(source_repo, target_repo, criteria)
	if err != nil {
		log.Fatal(err)
	}

	syncTaskId := associateCallReport.SpawnedTasks[0].TaskId
	fmt.Printf("TaskId: %v\n", syncTaskId)

	state := "init"
	for (state != "finished") && (state != "error") {
		task, _, terr := client.Tasks.GetTask(syncTaskId)
		fmt.Printf("task: %v\n", task)

		if task.State == "error" {
			fmt.Printf("error: %v\n", task.Error.Description)
		}

		if task.State == "finished" {
			for _, resultUnit := range task.Result.ResultUnits {
				fmt.Printf("rpm: %s\n", resultUnit.UnitKey)
			}
		}
		if task.State == "running" {
			fmt.Printf("----- progress --------\n")
			fmt.Printf("state: %v\n", task.State)
			fmt.Printf("progressReport: %v\n", task.ProgressReport)
		}
		state = task.State
		time.Sleep(10 * time.Millisecond)
		if terr != nil {
			log.Fatal(terr)
		}
	}
}
