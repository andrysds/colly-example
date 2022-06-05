package main

import (
	"log"
	"os"

	"github.com/andrysds/dropship-checker/checker"
	"github.com/andrysds/dropship-checker/csv"
	"github.com/andrysds/dropship-checker/partner"
	"github.com/subosito/gotenv"
)

const csvPathEnvKey = "CSV_PATH"

func main() {
	log.Println("starting...")

	gotenv.Load()

	csvPath := os.Getenv(csvPathEnvKey)
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalln("[ERROR] [opening csv file]", err)
	}

	r, err := csv.NewCSV(f)
	if err != nil {
		log.Fatalln("[ERROR] [NewCSV]", err)
	}

	p := partner.NewPartner()

	c := checker.NewChecker(r, p)

	if err := c.Check(); err != nil {
		log.Fatalln("[ERROR] [Check]", err)
	}

	log.Println("exiting...")
}
