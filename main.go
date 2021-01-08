package main

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"reflect" //testing data types. Remove when finished

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

	//assigning return of apiRequest() function to "response" variable
	response := apiRequest()
	
	fmt.Println("printing Json")
	fmt.Println(response)

	//calling ParseJSON function and passing response of apirequest() as a parameter
	parseJSON(response)
	//fmt.Println(response)

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

	//Create variable for the GET request and perform request with supplied variables from config file
	request, err := http.NewRequest("GET", "https://" + config.TenantURL + ".live.dynatrace.com/api/v2/problems", nil)
	//Request error handling
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//setting HTTP header for the GET request with supplied variables from config file
	request.Header.Set("Authorization", "Api-Token " + config.ApiToken)

	//Variable for the Response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()


	//variable to map response body 
	responseBody, err := ioutil.ReadAll(response.Body)

	//if statement to determine if the error variable is empty. If empty then log error 
	if err != nil  {
		log.Fatal(err)
		}
	//Printing Response body For testing. To be removed later.
	fmt.Println("Printing Response Body")
	fmt.Println(string(responseBody))
	//Add Error Handling for response body here.
	//Error handling should check for output of responseBody to see if token failed



	var responseObject Response

	//unmartial the responseData Json payload and assign to the responseObject variable
	json.Unmarshal(responseBody, &responseObject)
	//print the responseData variable as a string
	//fmt.Println(string(responseData))

	//returns the responseObject variable as a "Response" Type
	return responseObject

}

func parseJSON(jsonData Response){

	var totalProblems int 
	var infraProblems int = 0
	var serviceProblems int = 0
	//var topProblems []string
	totalProblems = jsonData.TotalCount


	//for loop to iterate over response data
	for i:= range jsonData.Problems{
		//fmt.Println(jsonData.Problems[i].ImpactLevel)
		//Checking data type
		//fmt.Println(reflect.TypeOf(jsonData.Problems[i].ImpactLevel))

		//if statement to check whether the problems "impactLevel" is SERVICES or INFRASTRUCTURE
		if jsonData.Problems[i].ImpactLevel == "SERVICES" {
			//Increment the "serviceProblems" variable by 1
			serviceProblems += 1
		} else if jsonData.Problems[i].ImpactLevel == "INFRASTRUCTURE" {
			//Increment the "infraProblems" variable by 1
			infraProblems += 1
		}

	}

	//testing output of variables
	fmt.Println("Printing total Problems")
	fmt.Println(totalProblems)
	fmt.Println("Printing Service Problems")
	fmt.Println(serviceProblems)
	fmt.Println("Printing Infrastructure Problems")
	fmt.Println(infraProblems)



}

