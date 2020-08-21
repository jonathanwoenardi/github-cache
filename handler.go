package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	client          *http.Client
	username        string
	password        string
	lastModifiedMap map[string]string
	commitDateMap   map[string]string
)

type APIResponse struct {
	CommitFullInfos []CommitFullInfo
}

type CommitFullInfo struct {
	Commit CommitInfo `json:"commit"`
}

type CommitInfo struct {
	Author Author `json:"author"`
}

type Author struct {
	Date string `json:"date"`
}

func getCommitDate(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	path := r.URL.Query().Get("path")
	if owner == "" || repo == "" {
		log.Print("owner and repo must be non empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := getURL(owner, repo, path)
	lastModified, _ := lastModifiedMap[url]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("NewRequest|err:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(username, password)
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do|err:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		response := &APIResponse{}
		err = json.NewDecoder(resp.Body).Decode(&response.CommitFullInfos)
		if err != nil {
			log.Printf("Decode|err:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseLastModified := resp.Header.Get("Last-Modified")
		lastModifiedMap[url] = responseLastModified
		var commitDate string
		if len(response.CommitFullInfos) != 0 {
			commitDate = response.CommitFullInfos[0].Commit.Author.Date
		}
		commitDateMap[url] = commitDate
		log.Printf("OK|url:%v|lastModified:%v|commitDate:%v", url, responseLastModified, commitDate)
		fmt.Fprintf(w, commitDate)
	case http.StatusNotModified:
		log.Printf("NotModified|url:%v|lastModified:%v|commitDate:%v", url, lastModified, commitDateMap[url])
		fmt.Fprintf(w, commitDateMap[url])
	default:
		log.Printf("StatusCode:%v", resp.StatusCode)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	// Debug rate limit
	rlLimit := resp.Header.Get("X-RateLimit-Limit")
	rlRemaining := resp.Header.Get("X-RateLimit-Remaining")
	rlReset := resp.Header.Get("X-RateLimit-Reset")
	log.Printf("RateLimit|limit:%v|remaining:%v|reset:%v", rlLimit, rlRemaining, rlReset)
}

func getURL(owner, repo, path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%v/%v/commits?path=%v", owner, repo, path)
}
