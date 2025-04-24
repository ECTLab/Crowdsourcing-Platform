package main

import (
	"crowdsourcing/pkg/app"
	"os"
	"time"

	inridereports "crowdsourcing/pkg/domains/in-ride-reports"

	log "github.com/sirupsen/logrus"
)

func main() {
	app.Init()

	err := inridereports.AggregateAllOnlineReports()
	if err != nil {
		log.WithError(err).Errorf("Error when running Online Report Aggregation job : %v", err)
		defer os.Exit(1)
	}

	// wait for prometheus to scrape metrics
	waitTime := 20
	log.Infof("Job is done. waiting %d seconds for prometheus scrapes to happen", waitTime)
	time.Sleep(time.Duration(waitTime * int(time.Second)))
}
