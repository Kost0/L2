package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

const defaultNTPServer = "pool.ntp.org"

func main() {
	now, err := ntp.Time(defaultNTPServer)
	if err != nil {
		log.New(os.Stderr, "", 0).Println(err)
		os.Exit(1)
	}
	fmt.Println(now.Format(time.RFC3339))
}
