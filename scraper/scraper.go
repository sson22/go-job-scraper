package scraper

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
var applyLinkURL string = "https://au.indeed.com/viewjob?jk="
//Scrape Indeed.com by keyword
func Scrape(keyword string){
	var baseURL string = "https://au.indeed.com/jobs?q="+keyword+"&limit=50"

	var totalJobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := countPages(baseURL)
	//Merge extracted jobs from each page
	for i :=0; i < totalPages; i++ {
		//Create 5 * Go channel between main() and getpages() to execute functions concurrently
		go getPage(i, baseURL, c)
	}

	for i := 0 ; i < totalPages; i++ {
		extractedJobs := <-c
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
func getPage(pageNumber int, baseURL string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make (chan extractedJob)
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
		//Create 50 * Go routines between getPages() and extractJob() to execute the function concurrently
		go extractJob(card, c)
		
	})
	for i:=0; i < searchCards.Length(); i++ {
		job := <- c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}
//Extract job information
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Find(".jcs-JobTitle").Attr("data-jk")
	title := CleanString(card.Find(".jcs-JobTitle>span").Text())
	location := CleanString(card.Find(".companyLocation").Text())
	company := CleanString(card.Find(".companyName").Text())
	salary := CleanString(card.Find(".salary-snippet-container>.attribute_snippet").Text())
	description := CleanString(card.Find(".job-snippet>ul>li").Text())
	c <- extractedJob{id:id, title:title, location:location, company:company, salary: salary, description:description}
}
//Get number of pages to inspect on URL using goquery
func countPages(baseURL string) int {
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
func CleanString(str string) string {
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

