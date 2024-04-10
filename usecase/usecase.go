package usecase

import (
	"OFFLINECOMICS/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CustomError struct {
	Message string
	Code    int
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func Fetch(URL string, n int) (utils.Result, error) {
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

var jobs = make(chan int, 100)
var results = make(chan utils.Result, 100)
var resultCollection []utils.Result
