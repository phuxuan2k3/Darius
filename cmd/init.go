package cmd

import (
	"net/http"
	"time"
)

func initMissfortuneClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return client
}
