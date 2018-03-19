build: build-overseer build-overlord

build-overseer:
	go build -o ./tmp/overseer ./overseer

build-overlord:
	go build -o ./tmp/overlord ./overlord

clean:
	rm tmp/*