package main

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	TotalCount int `json:"totalCount"`
	PageSize int `json:"pageSize"`
	//problems Problems `json:"problems"`
}


type Problems struct {
	ProblemId string `json:"problemId"`
	DisplayId string `json:"displayId"`
	Title string `json:"title"`
	ImpactLevel string `json:"impactLevel"`
	SeverityLevel string `json:"severityLevel"`
	Status string `json:"status"`
	AffectedEntities affectedEntities `json:"affectedEntities"`
}

type affectedEntities struct {
	EntityId entityId `json:"entityId"`
}

type entityId struct {

}


type Config struct {
	TenantURL string `json:"tenantURL"`
	ApiToken string `json:"Api-Token"`

}

func main() {
	fmt.Println("Welcome to the Top problem report")
	apiRequest()

}

// apiRequest used to reach to the Dynatrace API and pull back response data as a json
func apiRequest(){
	//Reading data from configuration file 
	configFile, err := ioutil.ReadFile("config.json")
	//error handling for configuration file
	if err != nil {
		fmt.Println(err)
	}

	//create variable for Configurations
	var config Config
	//unmarshal the json into the config variable
	json.Unmarshal(configFile, &config)

	//variable for the GET request
	response, err := http.Get("https://" + config.TenantURL + ".live.dynatrace.com/api/v2/problems?Api-Token=" + config.ApiToken)
	//if statement to determine if there is an error in the GET request. If there is an error then exit the program
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//variable to map response data 
	responseData, err := ioutil.ReadAll(response.Body)
	//if statement to determine if the error variable is empty. If empty then log error 
	if err != nil  {
		log.Fatal(err)
		}

	var responseObject Response


	json.Unmarshal(responseData, &responseObject)
	//print the responseData variable as a string
	fmt.Println(responseObject.TotalCount)
	fmt.Println(responseObject.PageSize)
	//fmt.Println(string(responseData))


}
