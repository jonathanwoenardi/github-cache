package main

import (
	"log"
	"net/http"
	"time"
)

const (
	configFile = "config.yml"
)

func main() {
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	client = NewHTTPClient(5*time.Second, 50)
	username = config.Username
	password = config.Password
	lastModifiedMap = make(map[string]string)
	commitDateMap = make(map[string]string)

	http.HandleFunc("/", getCommitDate)

	// Start HTTP server
	log.Printf("HTTP server started on %v", config.HTTPListen)
	err = http.ListenAndServe(config.HTTPListen, nil)
	if err != nil {
		log.Fatalf("ListenAndServe:%v", err)
	}
}
