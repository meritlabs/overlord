package messaging

import "net/http"

func DoPing() {
	http.Get("http://localhost:8080/ping")
}

func DoCheck() {
	http.Get("http://localhost:8080/check")
}
