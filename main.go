package main

import(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	//"reflect" //testing data types. Remove when finished
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

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
	ManagementZones []struct{
		Name string `json:"name"`
	} `json:"managementZones"`
}

type kv struct {
	Key   string
	Value int
}





func main() {
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

	fmt.Println("Welcome to the Top problem report")
	var managementZoneNames[]string
	items := make([]opts.BarData, 0)






	//Main for loop to iterate through get requests and parse json file
	for mz := range config.ManagementZones{
		//Print text to show which management zone we are querying
		fmt.Println("Querying Dynatrace instance for Management Zone :" , config.ManagementZones[mz].Name)
		//appending management zone name to the managementZoneNames array
		managementZoneNames = append(managementZoneNames, config.ManagementZones[mz].Name)

		//response variable for return value of this iteration of api request
		response := apiRequest(config.TenantURL, config.ApiToken, config.ManagementZones[mz].Name )
		//Variable for total count of problems in the JSON file
		totalProblemCount := returnTotalProblems(response)
		if totalProblemCount >=	500{
			fmt.Println("The problem payload is " ,  totalProblemCount , "Please reduce the number of problems to parse." )
			os.Exit(1)
		}
		//variable for count of infrastructure problems in JSON file
		infraProblemCount := returnInfraProblems(response)
		//variable for count of service problems in JSON file
		serviceProblemCount := returnServiceProblems(response)

		//variable for map of all problems in JSON file
		problemList := returnProblemList(response)


		//Testing output of functions here
		fmt.Println("Printing total number of problems")
		fmt.Println(totalProblemCount)
		fmt.Println("Printing number of infrastructure problems")
		fmt.Println(infraProblemCount)
		fmt.Println("Printing number of Service problems")
		fmt.Println(serviceProblemCount)
		fmt.Println("Printing Map of problems")
		fmt.Println(problemList)
		items = append(items, opts.BarData{Value: totalProblemCount})



	}

	setBarChart(managementZoneNames, items)






}

// apiRequest used to reach to the Dynatrace API and pull back response data as a json
func apiRequest(tenantURL string, apiToken string, managementZone string)(Response){

	//Create variable for the GET request and perform request with supplied variables from config file
	request, err := http.NewRequest("GET", "https://" + tenantURL + ".live.dynatrace.com/api/v2/problems?problemSelector=managementZones%28%22" + url.PathEscape(managementZone) + "%22%29", nil)
	//Request error handling
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	//setting HTTP header for the GET request with supplied variables from config file
	request.Header.Set("Authorization", "Api-Token " + apiToken)

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

	//Add Error Handling for response body here.
	//Error handling should check for output of responseBody to see if token failed

	var responseObject Response

	//unmartial the responseData Json payload and assign to the responseObject variable
	json.Unmarshal(responseBody, &responseObject)

	//returns the responseObject variable as a "Response" Type
	return responseObject

}



//function to take the Response json and output an integer count of total amount of problems
func returnTotalProblems(jsonData Response)(int){
	var totalProblems int 

	totalProblems = jsonData.TotalCount

	return totalProblems
}

//function to take the Response json and output an integer count of infrastructure problems
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

//function to take the Response json and output an integer count of Service problems
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


//function to take the Response json and output a map of the problems
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

    			problemList[key] = val + 1


			}else{
				key:= jsonData.Problems[i].AffectedEntities[y].Name
				problemList[key] = 1
		}
		}
	}
	sortProblemList(problemList)

	return problemList
}

func sortProblemList(problemList map[string]int)([]kv){
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


func setBarChart(managementZoneNames []string, problemCount []opts.BarData) {

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Problems sorted by Management Zones",
		Subtitle: "Default is set to past 2 hours",
	}))

	// Put data into instance
	bar.SetXAxis(managementZoneNames).
		AddSeries("Category A", problemCount)
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)

}





