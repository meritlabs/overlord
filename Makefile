build: build-overseer build-overlord

build-overseer:
	go build -o ./bin/overseer ./cmd/overseer/main.go

build-overlord:
	go build -o ./bin/overlord ./cmd/overlord/main.go

clean:
	rm overseer overlord