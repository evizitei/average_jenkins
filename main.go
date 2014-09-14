package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"code.google.com/p/gcfg"
)

type Config struct {
	Jenkins struct {
		Builds []string
		Host   string
	}
}

func (conf *Config) Urls() []string {
	urls := make([]string, 3, 10)
	for i, build := range conf.Jenkins.Builds {
		urls[i] = "http://" + conf.Jenkins.Host + "/job/" + build + "/api/json?tree=builds[number,id,timestamp,result,duration]"
	}
	return urls
}

type Project struct {
	Builds []struct {
		Duration  int    `json:"duration"`
		Id        string `json:"id"`
		Number    int    `json:"number"`
		Result    string `json:"result"`
		Timestamp int64  `json:"timestamp"`
	}
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

func main() {
	config, cnfErr := LoadConfig()
	if cnfErr != nil {
		fmt.Println("Configuration Error:", cnfErr)
		return
	}

	averages := make([]int, 3, 10)
	for i, url := range config.Urls() {
		averages[i] = AnalyzeBuild(url)
	}
	largest := selectMax(averages)
	mins := largest / 60000
	fmt.Println("Average interval:", mins, "minutes")
}

func selectMax(values []int) int {
	max := 0
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
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

func LoadConfig() (Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, "avgjenkins.gcfg")
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
