package main

import "code.google.com/p/gcfg"

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

func LoadConfig() (Config, error) {
	var config Config
	err := gcfg.ReadFileInto(&config, "avgjenkins.gcfg")
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
