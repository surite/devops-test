package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

type ProjectStatusRes struct{
	ProjectStatus Status `json:"projectStatus"`
}

type Status struct {
	Status              string	`json:"status"`
	Conditions	[]Condition	`json:"conditions"`
}

type Condition struct {
	Status	string	`json:"status"`
	MetricKey string `json:"metricKey"`
	Comparator string `json:"comparator"`
	ErrorThreshold string `json:"errorThreshold"`
	ActualValue	string `json:actualValue`
}

var(
	hostName string
	queryProjectKeys string
	queryBranchName string
)

func main() {

	flag.StringVar(&queryProjectKeys,"ProjectKeys","","Query Project Keys,Separated by commas, etc projectA,projectB")
	flag.StringVar(&queryBranchName,"BranchName","","Branch Name,etc master")
	flag.StringVar(&hostName,"HostName","sonarcloud.io","SonarQube host name")

	flag.Parse()

	projectKeys := strings.Split(queryProjectKeys,",")

	for _,projectKey := range projectKeys{

		queryProjectKey := projectKey

		if(queryBranchName != ""){
			queryProjectKey = projectKey + ":" + queryBranchName
		}

		checkProjectQualtiyStatus(queryProjectKey)
	}

	fmt.Println("check projects qualitygate status successfully!")
}

func checkProjectQualtiyStatus(projectKey string){
	apiURL := fmt.Sprintf("https://sonarcloud.io/api/qualitygates/project_status?projectKey=%s", projectKey)

	req, _ := http.NewRequest("GET", apiURL, nil)
	//set Authority Header
	setReqHeader(req)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var respBody = &ProjectStatusRes{}
	fmt.Println(res)
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		panic(fmt.Sprintf("ERROR: Failed to deserialize request body (%s)", err.Error()))
	}

	if respBody.ProjectStatus.Status != "OK" {
		panic(fmt.Sprintf("ERROR: %s project qualitygate status is not OK",projectKey))
	}

}

func setReqHeader(req *http.Request) {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", hostName)
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}
