package main

import (
	"OFFLINECOMICS/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const Url = "https://xkcd.com"

func fetch(n int) (*utils.Result, error) {

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// concatenate strings to get url; ex: https://xkcd.com/571/info.0.json
	url := strings.Join([]string{Url, fmt.Sprintf("%d", n), "info.0.json"}, "/")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("http request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http err: %v", err)
	}

	var data utils.Result

	// error from web service, empty struct to avoid disruption of process
	if resp.StatusCode != http.StatusOK {
		data = utils.Result{}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("json err: %v", err)
		}
	}

	resp.Body.Close()

	return &data, nil
}

type Job struct {
	number int
}

var jobs = make(chan Job, 100)
var results = make(chan utils.Result, 100)
var resultCollection []utils.Result

func allocateJobs(noOfJobs int) {
	for i := 0; i <= noOfJobs; i++ {
		jobs <- Job{i + 1}
	}
	close(jobs)
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		result, err := fetch(job.number)
		if err != nil {
			log.Printf("error in fetching: %v\n", err)
		}
		results <- *result
	}
	wg.Done()
}

func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i <= noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}

func getResults(done chan bool) {
	for result := range results {
		if result.Num != 0 {
			fmt.Printf("Retrieving issue #%d\n", result.Num)
			resultCollection = append(resultCollection, result)
		}
	}
	done <- true
}

func main() {
	// allocate jobs
	noOfJobs := 3000
	go allocateJobs(noOfJobs)

	// get results
	done := make(chan bool)
	go getResults(done)

	// create worker pool
	noOfWorkers := 100
	createWorkerPool(noOfWorkers)

	// wait for all results to be collected
	<-done

	// convert result collection to JSON
	data, err := json.MarshalIndent(resultCollection, "", "    ")
	if err != nil {
		log.Fatal("json err: ", err)
	}

	// write json data to file
	err = writeToFile(data)
	if err != nil {
		log.Fatal(err)
	}
}

func writeToFile(data []byte) error {
	f, err := os.Create("xkcd.json")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
