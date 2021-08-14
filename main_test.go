package main

import (
	"fmt"
	"github.com/ebarkie/aprs"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestProcessMessages(t *testing.T) {
	f := aprs.Frame{
		Src:  aprs.Addr{Call: "ZS1JPM", SSID: 15, Repeated: false},
		Dst:  aprs.Addr{Call: "APRS"},
		Path: aprs.Path{aprs.Addr{Call: "TCPIP", Repeated: true}},
		Text: ">Test",
	}
	password, err := strconv.Atoi(os.Getenv("APRS_PASSWORD"))
	if err != nil {
		t.Fatalf("could not parse APRS_PASSWORD")
	}
	err = f.SendIS(os.Getenv("APRS_SERVER"), aprs.Addr{Call: os.Getenv("APRS_USER")}, password)
	if err != nil {
		log.Printf("Upload error: %s", err)
	}

	fmt.Printf("%s\r\n", f)
}

func TestReconnectingSocket(t *testing.T) {
	hostInfo := os.Getenv("APRS_SERVER")
	call := aprs.Addr{Call: os.Getenv("APRS_USER")}
	password, err := strconv.Atoi(os.Getenv("APRS_PASSWORD"))
	if err != nil {
		t.Fatalf("could not parse APRS_PASSWORD")
	}

	frames := make(chan aprs.Frame)
	go aprs.SendTCPFromChannel(hostInfo, call, password, frames)

	f := aprs.Frame{
		Src:  aprs.Addr{Call: "", SSID: 15, Repeated: false},
		Dst:  aprs.Addr{Call: "APRS"},
		Path: aprs.Path{aprs.Addr{Call: "TCPIP", Repeated: true}},
		Text: ">Test",
	}

	frames<-f

	forever := make(chan bool)
	<-forever
}