package main

import (
	"fmt"
	"time"

	"com.github/meritlabs/overlord/overlord/messaging"
)

func main() {
	fmt.Printf("Overlord is ready.\n")

	for t := range time.NewTicker(30 * time.Second).C {
		fmt.Printf("Performing checks %v \n", t)
		messaging.DoPing()
		messaging.DoCheck()
	}
}
