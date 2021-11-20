all:
	gofmt -s -w .
	go build -o org-stats
run:
	go run main.go
