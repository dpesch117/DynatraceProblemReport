package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	//"reflect" //testing data types. Remove when finished
	"time"

	"github.com/dpesch117/DynatraceProblemReport/chart"
	"github.com/dpesch117/DynatraceProblemReport/problems"
	"github.com/go-echarts/go-echarts/v2/components"
)

//Config struct is used to map the config.JSON file located with this application
type Config struct {
	TenantURL       string `json:"tenantURL"`
	ApiToken        string `json:"Api-Token"`
	ManagementZones []struct {
		Name string `json:"name"`
	} `json:"managementZones"`
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
	//sankeyItems := make([]opts.SankeyNode, 0)
	problemData := make(map[string][]string)

	//Main for loop to iterate through get requests and parse json file
	for mz := range config.ManagementZones {
		//Print text to show which management zone we are querying
		fmt.Println("Querying Dynatrace instance for Management Zone :", config.ManagementZones[mz].Name)

		//response variable for return value of this iteration of api request
		response := problems.ApiRequest(config.TenantURL, config.ApiToken, config.ManagementZones[mz].Name)
		time.Sleep(2 * time.Second)
		//Variable for total count of problems in the JSON file
		totalProblemCount := problems.ReturnTotalProblems(response)
		if totalProblemCount >= 500 {
			fmt.Println("The problem payload is ", totalProblemCount, "Please reduce the number of problems to parse.")
			os.Exit(1)
		}

		//variable for count of infrastructure problems in JSON file
		infraProblemCount := problems.ReturnInfraProblems(response)
		//variable for count of service problems in JSON file
		serviceProblemCount := problems.ReturnServiceProblems(response)

		//variable for map of all problems in JSON file
		problemList := problems.ReturnProblemList(response)

		//Testing output of functions here
		fmt.Println("Printing total number of problems")
		fmt.Println(totalProblemCount)
		fmt.Println("Printing number of infrastructure problems")
		fmt.Println(infraProblemCount)
		fmt.Println("Printing number of Service problems")
		fmt.Println(serviceProblemCount)
		fmt.Println("Printing Map of problems")
		fmt.Println(problemList)

		//appending data to the problemData MAP. Current map looks like {Key: {value, value, value }}
		problemData[config.ManagementZones[mz].Name] = append(problemData[config.ManagementZones[mz].Name], strconv.Itoa(totalProblemCount))
		problemData[config.ManagementZones[mz].Name] = append(problemData[config.ManagementZones[mz].Name], strconv.Itoa(infraProblemCount))
		problemData[config.ManagementZones[mz].Name] = append(problemData[config.ManagementZones[mz].Name], strconv.Itoa(serviceProblemCount))
		fmt.Println("Printing problemData variable: ", problemData)
	}

	chart.BarChart(problemData)
	//fmt.Println(problemData)
	chart.Sankey(problemData)

	fmt.Println("running Dashboard() function")

	//Testing Adding multiple charts to a single dashboard
	page := components.NewPage()
	page.AddCharts(
		chart.GraphSankey(),
		chart.BarChart(problemData),
	)
	f, err := os.Create("test.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

}
