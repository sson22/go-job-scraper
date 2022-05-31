package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/sson22/job-scrapper/scraper"
)

const fileName string = "jobs.csv"
func handleHome(c echo.Context) error {
	return c.File("home.html")
}
func handleScrape(c echo.Context) error{
	defer os.Remove(fileName)
	keywords := strings.ToLower(scraper.CleanString(c.FormValue("keywords")))
	scraper.Scrape(keywords)
	return c.Attachment(fileName,fileName)
}
func main(){
 
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	godotenv.Load()
	port:=os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))


	
}