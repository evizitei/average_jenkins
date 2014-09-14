package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type Build struct {
	Duration  int    `json:"duration"`
	Id        string `json:"id"`
	Number    int    `json:"number"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
}

type Project struct {
	Builds []Build `json:"builds"`
}

func (proj *Project) averageDuration() int {
	summedDuration := 0
	count := 0
	for _, build := range proj.Builds {
		if build.Result == "SUCCESS" {
			summedDuration += build.Duration
			count += 1
		}
	}
	return int(math.Ceil(float64(summedDuration / count)))
}

func AnalyzeBuild(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in http request:", err)
		return -1
	}
	defer resp.Body.Close()
	var project Project
	jsonErr := json.NewDecoder(resp.Body).Decode(&project)

	if jsonErr != nil {
		fmt.Println("Error in JSON decoding:", jsonErr)
		return -1
	}
	return project.averageDuration()
}
