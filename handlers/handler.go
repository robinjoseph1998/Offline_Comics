package handlers

import (
	"OFFLINECOMICS/usecase"
	"OFFLINECOMICS/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReciveComicUrl(c *gin.Context) {
	var resultCollection []utils.Result
	var URL string
	if err := c.ShouldBind(&URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request", "error": err.Error()})
		return
	}
	no := 3000
	go usecase.AllocateJobs(no)

	done := make(chan bool)
	go usecase.GetResults(done, resultCollection)

	noOfworkers := 100
	go usecase.CreateWorkerPool(noOfworkers)
	<-done

	data, err := json.MarshalIndent(resultCollection, "", "    ")
	if err != nil {
		log.Fatal("json error: ", err)
	}
}
