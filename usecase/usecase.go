package usecase

import (
	"OFFLINECOMICS/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CustomError struct {
	Message string
	Code    int
}

const URL = "https://xkcd.com"

func (e CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func Fetch(n int) (utils.Result, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	url := strings.Join([]string{URL, fmt.Sprintf("%d", n), "info.0.json"}, "/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return utils.Result{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return utils.Result{}, err
	}
	var data utils.Result
	if resp.StatusCode != http.StatusOK {
		err := CustomError{
			Message: "Something Error",
			Code:    500,
		}
		return utils.Result{}, err
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return utils.Result{}, err
		}
	}
	resp.Body.Close()
	return data, nil

}

type Job struct {
	number int
}

var jobs = make(chan Job, 100)
var results = make(chan utils.Result, 100)
var resultCollection []utils.Result

func AllocateJobs(intJobs int) {
	for i := 0; i < intJobs; i++ {
		jobs <- Job{i + 1}
	}
	close(jobs)
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		result, err := Fetch(job.number)
		if err != nil {
			log.Fatalf("err in fetching", err)
			return
		}
		results <- *&result
	}
	wg.Done()
}
