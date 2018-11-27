all: tech-debt

.PHONY: client clean

tech-debt: client server
	rice append --exec server
	mv server tech-debt

client: client/src/main.js
	cd client && npm run build

server: server.go
	go build -o server

clean:
	rm -rf server tech-debt client/dist