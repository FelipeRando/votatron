package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordVotes(votingID *string, alternativeID *string) {
	go func() {
		tries := 0
		for {
			resp := vote(*votingID, *alternativeID)
			if tries == 10 {
				log.Fatal("Failed 10 times")
			}
			if resp.StatusCode != 200 {
				log.Printf("Vote Failed! Status: %v\n", resp.Status)
				tries++
			} else {
				log.Println("Voted successfully!")
				votesProcessed.Inc()
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

var (
	votesProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "votatron_processed_votes_total",
		Help: "The total number of processed votes",
	})
)

func vote(votingID string, alternativeID string) *http.Response {
	client := &http.Client{}
	data := url.Values{}
	data.Set("voting_id", votingID)
	data.Set("alternative_id", alternativeID)

	req, err := http.NewRequest("POST", "https://voting-vote-producer.r7.com/vote", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("authority", "voting-vote-producer.r7.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Mobile Safari/537.36")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://afazenda.r7.com")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://afazenda.r7.com/a-fazenda-12")
	req.Header.Add("accept-language", "en-US,en;q=0.9,pt-BR;q=0.8,pt;q=0.7")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return resp
}

var votingID *string
var alternativeID *string

func init() {
	votingID = flag.String("votingID", "", "The voting ID")
	alternativeID = flag.String("alternativeID", "", "The alternative ID")

	flag.Parse()
}
func main() {
	recordVotes(votingID, alternativeID)
	log.Printf("Started now! (Unix Time): %v\n", time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
	log.Println("Started http server for prometheus!")
}
