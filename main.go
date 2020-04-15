package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func shouldDownloadNewJobs() bool {
	info, err := os.Stat("descriptions.json")
	if os.IsNotExist(err) {
		log.Println("Jobs file does not exist.")
		return true
	}

	currentYear, currentMonth, currentDay := time.Now().Date()
	fileYear, fileMonth, fileDay := info.ModTime().Date()
	if currentYear == fileYear && currentMonth == fileMonth && currentDay == fileDay {

		log.Println("Jobs file already downloaded today.")
		return false
	}

	log.Println("Jobs file is old.")
	return true
}

func writeToFile(jobDescriptions []TEJob) error {
	bytes, err := json.Marshal(jobDescriptions)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("descriptions.json", bytes, 0644)
}

func main() {

	var data TEData
	if shouldDownloadNewJobs() {
		log.Println("Downloading new jobs...")
		data = FetchJobs()

		var jobDescriptions []TEJob

		for _, doc := range data.Response.Docs {
			descriptions := FetchJobDescription(doc.ID)
			jobDescriptions = append(jobDescriptions, descriptions...)
		}

		writeToFile(jobDescriptions)
		log.Println("Jobs file saved.")
	}
}
