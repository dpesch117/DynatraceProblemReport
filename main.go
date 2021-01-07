package main

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"//reflect can be removed. It is purely for testing purposes
)

type Response struct {
	TotalCount int `json:"totalCount"`
	PageSize int `json:"pageSize"`
	Problems []struct {
		ProblemId string `json:"problemId"`
		DisplayId string `json:"displayId"`
		Title string `json:"title"`
		ImpactLevel string `json:"impactLevel"`
		SeverityLevel string `json:"severityLevel"`
		Status string `json:"status"`
		AffectedEntities []struct{
			EntityId struct {
				Id string `json:"id"`
				Type string `json:"type"`
			} `json:"entityId"`
		} `json:"affectedEntities"`
		ImpactedEntities []struct{
			EntityId struct {
				Id string `json:"id"`
				Type string `json:"type"`
			} `json:"entityId"`
			Name string `json:"name"`
		} `json:"impactedEntities"`
		RootCauseEntity string `json:"rootCauseEntity"`
		ManagementZones []struct{
			Id string `json:"id"`
			Name string `json:"name"`
		} `json:"managementZones"`
		EntityTags []struct{
			Context string `json:"context"`
			Key string `json:"key"`
			Value string `json:"value"`
			StringRepresentation string `json:"stringRepresentation"`
		} `json:"entityTags"`
		ProblemFilters []struct{
			Id string `json:"id"`
			Name string `json:"name"`
		} `json:"problemFilters"`
		StartTime int `json:"startTime"`
		Endtime int `json:"endTime"`
	} `json:"problems"`
}

type Config struct {
	TenantURL string `json:"tenantURL"`
	ApiToken string `json:"Api-Token"`
}

func main() {
	fmt.Println("Welcome to the Top problem report")
	response := apiRequest()
	fmt.Println("printing Json")
	fmt.Println(reflect.TypeOf(response))
	fmt.Println(response)

}

// apiRequest used to reach to the Dynatrace API and pull back response data as a json
func apiRequest()(Response){
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
//testing for loop to iterate over response data
/*
	for i:= range responseObject.Problems{

		fmt.Println(responseObject.Problems[i].ManagementZones)

	}
*/

	//unmartial the responseData Json payload and assign to the responseObject variable
	json.Unmarshal(responseData, &responseObject)
	//print the responseData variable as a string
	//fmt.Println(string(responseData))

	//returns the responseObject variable as a "Response" Type
	return responseObject

}

