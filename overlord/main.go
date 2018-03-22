package main

import (
	"fmt"
	"time"

	"github.com/meritlabs/overlord/overlord/messaging"
)

var (
	pingDelay  = 30
	checkDelay = 300 // 5 minutes
	ticks      = checkDelay / pingDelay
)

func main() {
	fmt.Printf("Overlord is ready.\n")

	countdown := ticks
	for t := range time.NewTicker(time.Duration(pingDelay) * time.Second).C {
		fmt.Printf("Performing checks %v, %v\n", t, countdown)
		countdown--

		go messaging.DoPing()

		if countdown == 0 {
			go messaging.DoCheck()
			countdown = ticks
		}
	}
}
