package main

import (
	"bufio"
	"fmt"
	log2 "log"
	"os"
	"strconv"
	"sync"

	"github.com/toorop/go-ovh-logs"
)

// conf
var level uint8
var proto ovhlogs.Protocol
var endpoint string
var token string

var logs *ovhlogs.OvhLogs
var wg *sync.WaitGroup

func log(line string) {
	entry := ovhlogs.Entry{
		FullMessage: line,
		Level:       level,
	}
	logs.Send(entry)
	wg.Done()
}

func main() {
	var err error
	// get conf from ENV

	// level
	envLevel := os.Getenv("OVHLOGGER_LOGLEVEL")
	if envLevel == "" {
		log2.Fatal("OVHLOGGER_LOGLEVEL not found in ENV")
	}
	l, err := strconv.ParseUint(envLevel, 10, 64)
	if err != nil {
		log2.Fatal(err)
	}
	level = uint8(l)

	// proto
	envProto := os.Getenv("OVHLOGGER_PROTO")
	switch envProto {
	case "udp":
		proto = ovhlogs.GelfUDP
	case "tcp":
		proto = ovhlogs.GelfTCP
	case "tls":
		proto = ovhlogs.GelfTLS
	default:
		log2.Fatal(fmt.Sprintf("OVHLOGGER_PROTO not supported: %s", envProto))
	}

	// endpoint
	endpoint = os.Getenv("OVHLOGGER_ENDPOINT")
	if envLevel == "" {
		log2.Fatal("OVHLOGGER_ENDPOINT not found in ENV")
	}

	// token
	token = os.Getenv("OVHLOGGER_TOKEN")
	if envLevel == "" {
		log2.Fatal("OVHLOGGER_TOKEN not found in ENV")
	}

	// wait group
	wg = new(sync.WaitGroup)

	// init logger
	logs = ovhlogs.New(endpoint, token, proto, ovhlogs.CompressNone, false)

	// read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		go log(scanner.Text())
		fmt.Println(scanner.Text())
	}
	wg.Wait()
}
