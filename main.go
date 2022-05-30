package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
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
var applyLinkURL string = "https://au.indeed.com/viewjob?jk="

func main(){
	var totalJobs []extractedJob
	totalPages := countPages()
	//Merge extracted jobs from each page
	for i :=0; i < totalPages; i++ {
		extractedJobs :=getPage(i)
		totalJobs = append(totalJobs, extractedJobs...)
	}
	writeJobs(totalJobs)

	fmt.Println("Job extraction done for total", len(totalJobs),"jobs" )

}

//Wrtie jobs on file
func writeJobs(jobs[] extractedJob){
	file, err := os.Create("jobs.csv")
	checkError(err)
	w:= csv.NewWriter(file)
	//Flush data on the file
	defer w.Flush()

	headers := []string{"LINK","Title","Location","Company","Salary","Description"}
	wErr := w.Write(headers)
	checkError(wErr) 

	for _, job := range jobs {
		jobSlice := []string{applyLinkURL + job.id, job.title, job.location, job.company, job.salary, job.description}
		jobWErr := w.Write(jobSlice)
		checkError(jobWErr)
	}
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

