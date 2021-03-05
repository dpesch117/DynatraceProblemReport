# DynatraceProblemReport

**About:**
This Dynatrace Problem Report project is an attempt to make actionable data off of the Dynatrace platform "problems". When there are many problems being opened it can become overwhelming to find a point where you need to begin. The idea behind is to pull Data from the Dynatrace V2 API's, parse and visualize in a user friendly way.

Problem Reference - https://www.dynatrace.com/support/help/how-to-use-dynatrace/problem-detection-and-analysis/

I originally created a (closed source) problem report for one of my clients. The intent of this project was to rewrite the report in Go, utilize the newly deployed Dynatrace API v2 endpoints and open source it so anyone can run it on their Dynatrace cluster.

**Requirements:**
Go 1.15


**Instructions:**


**Download the Source code above.**

**Create the go.MOD file**
```
$go mod init https://github.com/dpesch117/DynatraceProblemReport/
```
**Create a config.json file with the following format:**
```
{
	"tenantURL":"", //Environment ID here
	"Api-Token": "", //API Token with API V2 "Read Problems" Permissions
	"managementZones" : [
		{
			"name" : ""
		}
	]
}
```
**Run the go progam**
```
$go run main.go
```

At this point you will find output HTML files inside your program's main directory. These files will be changed & updated as the project gets further along. 
