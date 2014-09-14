package main

import "fmt"

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
