package main

import (
	"diskstatus/server"
	"diskstatus/service"
	"diskstatus/unit"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// args
var (
	interval        = flag.Duration("interval", time.Second, "the timestamp interval disk info will be retrieved.")
	address         = flag.String("addr", "localhost:8080", "the adress server will be listening on.")
	shouldLogToFile = flag.Bool("log", false, "enables logging the retrieved data to a diskstatus.log file.")
	byteUnit        = flag.String("byteunit", "GB", "The byte unit to be used as storage metrics. Options are: B,KB,MB,GB,TB")
)


var (
	// Server
	svr *http.Server = nil
)

func main() {
	flag.Parse()

	go graceFullShutdown()

	pr, pw := io.Pipe()
	defer pw.Close()
	defer pr.Close()

	unitByte, err := filterByteUnitOption(*byteUnit)

	if err != nil {
		log.Fatal(err)
	}

	diskStatusService := service.New(*interval, unitByte)

	// our writer
	var m io.Writer
	if *shouldLogToFile {
		l := createLogFile()
		defer l.Close()
		m = io.MultiWriter(pw, l)
	} else {
		m = pw
	}

	buf := make([]byte, 7)

	go func() {
		for {
			_, err := pr.Read(buf)

			if err != nil {
				println(err.Error())
				break
			}
		}
	}()

	go diskStatusService.Execute(m)

	svr = server.CreateNewServer(*address, os.Stderr, server.TextStreamHTTPHandler(buf, *interval))
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func createLogFile() *os.File {
	f, err := os.Create("diskstatus.log")

	if err != nil {
		log.Fatal(err)
	}

	return f
}

func filterByteUnitOption(unitString string) (unit.ByteUnit, error) {
	switch unitString {
	case "B":
		return unit.Byte, nil
	case "KB":
		return unit.KByte, nil
	case "MB":
		return unit.MByte, nil
	case "GB":
		return unit.GByte, nil
	case "TB":
		return unit.TByte, nil
	default:
		return unit.ByteUnit(-1), errors.New("invalid byte unit")
	}
}

func graceFullShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sign := <- c
	println("SIGNAL: ", sign.String())

	if svr != nil {
		svr.Close()
	}
	os.Exit(0)
}