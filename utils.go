package main

import (
	"net"
	"net/http"
	"time"
)

func NewHTTPClient(timeout time.Duration, maxHosts int) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxHosts,
			MaxConnsPerHost:     maxHosts,
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 60 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: timeout,
		},
	}
}
