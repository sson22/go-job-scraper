package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://au.indeed.com/jobs?q=software&limit=50"

func main(){
	totalPages := countPages()
	fmt.Println(totalPages)

	for i :=0; i < totalPages; i++ {
		getPage(i)
	}

}
func getPage(pageNumber int) {
	pageURL := baseURL + "&start="+ strconv.Itoa(pageNumber*50)
	fmt.Println(pageURL)
}
//Get number of pages to inspect on URL using goquery
func countPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkError(err)
	checkStatusCode(res)
	//Close function to prevent memory leak
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err) 
	//Find div that has classname pagination from the URL
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection){
		//Count number of pages on the URL
		pages = s.Find("a").Length()
	})
	return pages
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