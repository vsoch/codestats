all:
	gofmt -s -w .
	go build -o codestats
run:
	go run main.go
