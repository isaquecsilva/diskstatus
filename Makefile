build:
	go build -ldflags='-s -w' -o bin/diskstatus.exe diskstatus.go
run:
	go run diskstatus.go -interval 1m
help:
	go run diskstatus.go -help