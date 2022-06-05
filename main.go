package main

import (
	"log"
	"os"

	"github.com/andrysds/colly-example/checker"
	"github.com/andrysds/colly-example/csv"
	"github.com/andrysds/colly-example/partner"
	"github.com/subosito/gotenv"
)

const csvPathEnvKey = "CSV_PATH"

func main() {
	log.Println("starting...")

	gotenv.Load()

	csvPath := os.Getenv(csvPathEnvKey)
	f, err := os.Open(csvPath)
	if err != nil {
		log.Println("[ERROR] [opening csv file]", err)
	}

	r, err := csv.NewCSV(f)
	if err != nil {
		log.Println("[ERROR] [NewCSV]", err)
	}

	p := partner.NewPartner()

	c := checker.NewChecker(r, p)

	if err := c.Check(); err != nil {
		log.Println("[ERROR] [Check]", err)
	}

	log.Println("exiting...")
}
