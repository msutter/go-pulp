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
	apiEndpoint := "pulp-lab-1.test"

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

	orFilter := pulp.NewFilter()
	orFilter.Operator = "$or"
	for _, fileName := range fileNames {
		// only match exact filename
		regex := fmt.Sprintf("^%s$", fileName)
		orFilter.AddExpression("filename", "$regex", regex)
	}

	subAndFilter := pulp.NewFilter()
	subAndFilter.Operator = "$and"
	subAndFilter.AddExpression("filename", "$regex", "node")

	subOrFilter := pulp.NewFilter()
	subOrFilter.Operator = "$or"
	subOrFilter.AddExpression("filename", "$regex", "nonono")

	SecondSubOrFilter := pulp.NewFilter()
	SecondSubOrFilter.Operator = "$or"
	SecondSubOrFilter.AddExpression("filename", "$regex", "wrong")

	orFilter.AddSubFilter(subAndFilter)
	orFilter.AddSubFilter(subOrFilter)
	orFilter.AddSubFilter(SecondSubOrFilter)

	criteria := pulp.NewUnitAssociationCriteria()
	criteria.AddFilter(orFilter)

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

		// find filenames for the results
		// resultFilter := pulp.NewFilter()
		if task.State == "finished" {
			fileNames, ubtrErr := client.Units.GetUnitFileNamesByTaskResult(task.Result)

			if ubtrErr != nil {
				fmt.Println(ubtrErr.Error())
				log.Fatal(ubtrErr)
			}
			for _, fileName := range fileNames {
				fmt.Printf("fileName: %v\n", fileName)
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
	fmt.Printf("-------------------------------------- SEARCH -------------------------------------------\n")

	name := "nodetree"
	regex := fmt.Sprintf("^%s$", name)

	searchCriteria := pulp.NewSearchCriteria()
	resultFields := []string{
		"name",
		"version",
		"arch",
		"release",
		"filename",
		"repository_memberships",
		"requires",
	}
	_ = resultFields
	// searchCriteria.AddFields(resultFields)

	// build query
	searchFilter := pulp.NewFilter()
	searchFilter.Operator = "$or"
	searchFilter.AddExpression("name", "$regex", regex)
	searchCriteria.AddFilter(searchFilter)

	units, _, uerr := client.Units.SearchUnits("rpm", searchCriteria)
	for _, unit := range units {
		fmt.Printf("unit: %v\n", unit)
	}

}
