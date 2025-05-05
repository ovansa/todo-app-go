package util

import (
	"log"
	"net/http"
	"time"
)

type SelfPinger struct {
	appURL         string
	pingInterval   time.Duration
	stopChan       chan struct{}
	isServerActive bool
}

func NewSelfPinger(appURL string, pingIntervalMinutes int) *SelfPinger {
	return &SelfPinger{
		appURL:         appURL,
		pingInterval:   time.Duration(pingIntervalMinutes) * time.Minute,
		stopChan:       make(chan struct{}),
		isServerActive: false,
	}
}

func (sp *SelfPinger) Start() {
	go func() {
		// Wait a bit for the server to fully initialize
		time.Sleep(10 * time.Second)
		sp.isServerActive = true

		ticker := time.NewTicker(sp.pingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sp.ping()
			case <-sp.stopChan:
				log.Println("Self-pinger stopped")
				return
			}
		}
	}()

	log.Printf("Self-pinger started. Will ping %s every %v", sp.appURL, sp.pingInterval)
}

func (sp *SelfPinger) Stop() {
	close(sp.stopChan)
}

// ping sends a request to the application URL
func (sp *SelfPinger) ping() {
	if !sp.isServerActive {
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(sp.appURL)
	if err != nil {
		log.Printf("Self-ping failed: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Self-ping performed. Response code: %d", resp.StatusCode)
}
