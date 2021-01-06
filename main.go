package main

import(
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Welcome to the Top problem report")
	apiRequest()

}

// apiRequest used to reach to the Dynatrace API and pull back response data as a json
func apiRequest(){
	//variable for the GET request
	response, err := http.Get("")
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
	//print the responseData variable as a string
	fmt.Println(string(responseData))


}