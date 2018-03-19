package messaging

import (
	"fmt"
	"net/http"
)

type Ping struct {
	ip       string
	isOnline bool
}

var (
	ips = []string{"127.0.0.1:8080", "192.168.1.6:8080"} // grab it from config
)

func DoPing() {
	pingChannel := make(chan Ping)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s/ping", ip)
			_, err := http.Get(host)

			if err != nil {
				pingChannel <- Ping{ip, false}
				return
			}

			pingChannel <- Ping{ip, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-pingChannel:
			if msg.isOnline {
				fmt.Printf("Node %s is online\n", msg.ip)
			} else {
				fmt.Printf("Node %s is offline\n", msg.ip)
			}
		}
	}
}

func DoCheck() {
	checkChannel := make(chan Ping)

	for _, ip := range ips {
		go func(ip string) {
			host := fmt.Sprintf("http://%s/check", ip)
			_, err := http.Get(host)

			if err != nil {
				checkChannel <- Ping{ip, true}
				return
			}

			checkChannel <- Ping{ip, true}
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		select {
		case msg := <-checkChannel:
			fmt.Printf("Node %s is checked\n", msg.ip)
		}
	}
}
