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

//Response struct is used to map the JSON file received in HTTP Response from Dynatrace Problem API v2
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
			Name string `json:"name"`
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
//Config struct is used to map the config.JSON file located with this application
type Config struct {
	TenantURL string `json:"tenantURL"`
	ApiToken string `json:"Api-Token"`
}


func main() {
	fmt.Println("Welcome to the Top problem report")

	//assigning return of apiRequest() function to "response" variable
	response := apiRequest()
	
	//calling ParseJSON function and passing response of apirequest() as a parameter
	parseJSON(response)


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
	//fmt.Println("Printing Response Body")
	//fmt.Println(string(responseBody))
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
	//Variable for total count of problems in the JSON file
	totalProblemCount := returnTotalProblems(jsonData)
	//variable for count of infrastructure problems in JSON file
	infraProblemCount := returnInfraProblems(jsonData)
	//variable for count of service problems in JSON file
	serviceProblemCount := returnServiceProblems(jsonData)

	//variable for map of all problems in JSON file
	problemList := returnProblemList(jsonData)

	fmt.Println("Printing total number of problems")
	fmt.Println(totalProblemCount)
	fmt.Println("Printing number of infrastructure problems")
	fmt.Println(infraProblemCount)
	fmt.Println("Printing number of Service problems")
	fmt.Println(serviceProblemCount)
	fmt.Println("Printing Map of problems")
	fmt.Println(problemList)


}

func returnTotalProblems(jsonData Response)(int){
	var totalProblems int 

	totalProblems = jsonData.TotalCount

	return totalProblems

}

func returnInfraProblems(jsonData Response)(int){

	var infraProblems int = 0

	for i:= range jsonData.Problems{


		//if statement to check whether the problems "impactLevel" is SERVICES or INFRASTRUCTURE
		if jsonData.Problems[i].ImpactLevel == "INFRASTRUCTURE" {
			//Increment the "serviceProblems" variable by 1
			infraProblems += 1
		}
	}

	return infraProblems

}

func returnServiceProblems(jsonData Response)(int){

	var serviceProblems int = 0

	for i:= range jsonData.Problems{


		//if statement to check whether the problems "impactLevel" is SERVICES or INFRASTRUCTURE
		if jsonData.Problems[i].ImpactLevel == "SERVICES" {
			//Increment the "serviceProblems" variable by 1
			serviceProblems += 1
		}
	}

	return serviceProblems
}

func returnProblemList(jsonData Response)(map[string]int){

	//created a list of problems based on the Problems struct
	problemList := make(map[string]int)


	//for loop to iterate over response data
	for i:= range jsonData.Problems{

		//for loop to iterate over the objects in the AffectedEntities of the JSON
		for y := range jsonData.Problems[i].AffectedEntities{

			//Assign variable 'key' to this iteration of the AffectedEntities.Name data
			key:= jsonData.Problems[i].AffectedEntities[y].Name


			//create an if statement to check if there is a value for this item
			if val, ok := problemList[key]; ok {
    			fmt.Println("problemList[key] is equal to :" , problemList[key] , " and val is equal to :" , val)
    			problemList[key] = problemList[key] + 1


			}else{
				key:= jsonData.Problems[i].AffectedEntities[y].Name
				problemList[key] = 1
		}
		}
	}

	return problemList

}



