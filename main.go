package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	company string
	salary string
	description string
}
var baseURL string = "https://au.indeed.com/jobs?q=software&limit=50"

func main(){
	var totalJobs []extractedJob
	totalPages := countPages()
	for i :=0; i < totalPages; i++ {
		extractedJobs :=getPage(i)
		totalJobs = append(totalJobs, extractedJobs...)
	}
	fmt.Println(totalJobs)

}

//Get each page and find each job-id
func getPage(pageNumber int)[]extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start="+ strconv.Itoa(pageNumber*50)
	res, err := http.Get(pageURL)
	checkError(err)
	checkStatusCode(res)
	//Close function to prevent memory leak
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	//Find id of each job and extract job information
	searchCards := doc.Find(".cardOutline")
	searchCards.Each(func(i int, card *goquery.Selection){
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}
//Extract job information
func extractJob(card *goquery.Selection) extractedJob{
	id, _ := card.Find(".jcs-JobTitle").Attr("data-jk")
	title := cleanString(card.Find(".jcs-JobTitle>span").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	company := cleanString(card.Find(".companyName").Text())
	salary := cleanString(card.Find(".salary-snippet-container>.attribute_snippet").Text())
	description := cleanString(card.Find(".job-snippet>ul>li").Text())
	return extractedJob{id:id, title:title, location:location, company:company, salary: salary, description:description}
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

//Trim empty string
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

