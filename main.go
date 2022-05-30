package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://au.indeed.com/jobs?q=software&limit=50"

func main(){
	pages := getPages()
	fmt.Println(pages)

}

//Get pages using goquery
func getPages() int {
	res, err := http.Get(baseURL)
	checkError(err)
	checkStatusCode(res)

	//Close function to prevent memory leak
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	//Find div that has classname pagination from the URL
	doc.Find(".pagination").Each()
	return 0
}

//Check error and log the error
func checkError(err error){
	if err != nil {
		log.Fatalln(err)
	}
}
//Check status code and log the code if request was not fulfilled
func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}