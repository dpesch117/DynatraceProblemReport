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

func apiRequest(){

	response, err := http.Get("")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil  {
		log.Fatal(err)
		}
	
	fmt.Println(string(responseData))


}