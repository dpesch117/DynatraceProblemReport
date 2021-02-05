package problems

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
)

type Response struct {
	TotalCount int `json:"totalCount"`
	PageSize   int `json:"pageSize"`
	Problems   []struct {
		ProblemId        string `json:"problemId"`
		DisplayId        string `json:"displayId"`
		Title            string `json:"title"`
		ImpactLevel      string `json:"impactLevel"`
		SeverityLevel    string `json:"severityLevel"`
		Status           string `json:"status"`
		AffectedEntities []struct {
			EntityId struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"entityId"`
			Name string `json:"name"`
		} `json:"affectedEntities"`
		ImpactedEntities []struct {
			EntityId struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"entityId"`
			Name string `json:"name"`
		} `json:"impactedEntities"`
		RootCauseEntity string `json:"rootCauseEntity"`
		ManagementZones []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"managementZones"`
		EntityTags []struct {
			Context              string `json:"context"`
			Key                  string `json:"key"`
			Value                string `json:"value"`
			StringRepresentation string `json:"stringRepresentation"`
		} `json:"entityTags"`
		ProblemFilters []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"problemFilters"`
		StartTime int `json:"startTime"`
		Endtime   int `json:"endTime"`
	} `json:"problems"`
}

type kv struct {
	Key   string
	Value int
}

func apiRequest(tenantURL string, apiToken string, managementZone string) Response {

	//Create variable for the GET request and perform request with supplied variables from config file
	request, err := http.NewRequest("GET", "https://"+tenantURL+".live.dynatrace.com/api/v2/problems?from=now-2d&problemSelector=managementZones%28%22"+url.PathEscape(managementZone)+"%22%29", nil)
	//Request error handling
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//setting HTTP header for the GET request with supplied variables from config file
	request.Header.Set("Authorization", "Api-Token "+apiToken)

	//Variable for the Response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//variable to map response body
	responseBody, err := ioutil.ReadAll(response.Body)

	//if statement to determine if the error variable is empty. If empty then log error
	if err != nil {
		log.Fatal(err)
	}

	//Add Error Handling for response body here.
	//Error handling should check for output of responseBody to see if token failed

	var responseObject Response

	//unmartial the responseData Json payload and assign to the responseObject variable
	json.Unmarshal(responseBody, &responseObject)

	//returns the responseObject variable as a "Response" Type
	return responseObject

}

//function to take the Response json and output an integer count of total amount of problems
func returnTotalProblems(jsonData Response) int {
	var totalProblems int

	totalProblems = jsonData.TotalCount

	return totalProblems
}

//function to take the Response json and output an integer count of infrastructure problems
func returnInfraProblems(jsonData Response) int {

	var infraProblems int = 0

	for i := range jsonData.Problems {

		//if statement to check whether the problems "impactLevel" is SERVICES or INFRASTRUCTURE
		if jsonData.Problems[i].ImpactLevel == "INFRASTRUCTURE" {
			//Increment the "serviceProblems" variable by 1
			infraProblems += 1
		}
	}

	return infraProblems

}

//function to take the Response json and output an integer count of Service problems
func returnServiceProblems(jsonData Response) int {

	var serviceProblems int = 0

	for i := range jsonData.Problems {

		//if statement to check whether the problems "impactLevel" is SERVICES or INFRASTRUCTURE
		if jsonData.Problems[i].ImpactLevel == "SERVICES" {
			//Increment the "serviceProblems" variable by 1
			serviceProblems += 1
		}
	}

	return serviceProblems
}

//function to take the Response json and output a map of the problems
func returnProblemList(jsonData Response) map[string]int {

	//created a list of problems based on the Problems struct
	problemList := make(map[string]int)

	//for loop to iterate over response data
	for i := range jsonData.Problems {

		//for loop to iterate over the objects in the AffectedEntities of the JSON
		for y := range jsonData.Problems[i].AffectedEntities {

			//Assign variable 'key' to this iteration of the AffectedEntities.Name data
			key := jsonData.Problems[i].AffectedEntities[y].Name

			//create an if statement to check if there is a value for this item
			if val, ok := problemList[key]; ok {

				problemList[key] = val + 1

			} else {
				key := jsonData.Problems[i].AffectedEntities[y].Name
				problemList[key] = 1
			}
		}
	}
	sortProblemList(problemList)

	return problemList
}

func sortProblemList(problemList map[string]int) []kv {
	var ss []kv

	for k, v := range problemList {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	//Commented out for loop that shows how to access keys and values in the []kv array
	//for _, kv := range ss {
	//
	//	fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	//}
	fmt.Println("Printing  SS")

	fmt.Println(ss)

	return ss
}
