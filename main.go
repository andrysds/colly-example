// https://github.com/andrysds/colly-example

package main

import (
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/subosito/gotenv"
)

func main() {
	log.Println("starting...")

	gotenv.Load()

	c := colly.NewCollector()

	c.Visit(os.Getenv("LOGIN_URL"))

	log.Println("exiting...")
}
