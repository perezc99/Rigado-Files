package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
)

const configFile = "config.toml"

var defaultEndpoint = "https://thesimpsonsquoteapi.glitch.me:443/quotes"

type Config struct {
	Endpoint string `toml:"endpoint"`
}
type SimpsonsQuote struct {
	Quote              string `json:"quote"`
	Character          string `json:"character"`
	Image              string `json:"image"`
	CharacterDirection string `json:"characterDirection"`
}

func loadConfig() (Config, error) {
	config := Config{}
	snapDir := os.Getenv("SNAP_DATA")
	data, err := os.ReadFile(filepath.Join(snapDir, configFile))
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(data, &config)
	return config, err
}

func queryAPI(endpoint string) {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("error querying API: %v", err)
		return
	}
	if resp.StatusCode > 399 {
		log.Fatalf("unexpected response %v - %v", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()

	var q []SimpsonsQuote

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		return
	}

	err = json.Unmarshal(body, &q)
	if err != nil {
		log.Printf("error reading quotes API - %v", err)
		return
	}
	quote := q[0]
	fmt.Printf("\"%s\" - %s\n", quote.Quote, quote.Character)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Printf("error loading config: %v", err)
		log.Printf("using default endpoint: %v", defaultEndpoint)
		config.Endpoint = defaultEndpoint
	}

	ticker := time.NewTicker(10 * time.Second)
	// Creating channel variable
	// Defaultly set as false
	channel_var := make(chan bool)

	// Using anonymous function to create a Goroutine
	// that holds the for-loop.
	go func() {
		queryAPI(config.Endpoint)
		for {
			select {
			// call queryAPI every 10 seconds
			case <-ticker.C:
				queryAPI(config.Endpoint)
			//  exit loop/function once sleep counter timesout
			case <- channel_var:
				return
			}
		}
	} ()

	// Setting up the number of loops;
	// queryAPI will be called the number of loops
	// specified plus one
	number_of_loops := 5
	
	// Sleep time is calculated by multiplying 
	// the number of loops with the number of seconds
	// set for the ticker cycle.  
	sleep_time := number_of_loops * 10

	// Start sleep function
	time.Sleep(sleep_time * time.Second)

	// Terminate ticker once sleep function timesout
	ticker.Stop()
	
	// Set channel variable to true to end Goroutine
	// holding the for loop.
	channel_var <- true

	// Print confirmation application has stopped running
	fmt.Printf("The application has completed!")

}