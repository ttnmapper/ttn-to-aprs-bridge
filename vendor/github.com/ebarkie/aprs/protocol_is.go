// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package aprs

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func genLogin(user Addr, pass int) string {
	return fmt.Sprintf("user %s pass %d vers %s %s", user, pass, SwName, SwVers)
}

func readLine(conn net.Conn, d time.Duration) (string, error) {
	if d > 0 {
		err := conn.SetReadDeadline(time.Now().Add(d))
		if err != nil {
			return "", err
		}
	} else {
		err := conn.SetReadDeadline(time.Time{})
		if err != nil {
			return "", err
		}
	}
	s, err := bufio.NewReader(conn).ReadString('\n')
	return strings.TrimSpace(s), err
}

// GenPass generates a verification passcode for the given station.
func GenPass(call string) (pass uint16) {
	// Refer to aprsc:
	// https://github.com/hessu/aprsc

	// Upper case callsign and strip SSID if it was included
	c := strings.ToUpper(call)
	dash := strings.Index(c, "-")
	if dash > -1 {
		c = c[:dash]
	}

	pass = 0x73e2 // The key/seed.
	for i := 0; i < len(c); i += 2 {
		pass ^= uint16(c[i]) << 8
		pass ^= uint16(c[i+1])
	}

	// Mask off the high bit so number is always positive
	pass &= 0x7fff
	return
}

// RecvIS receives APRS-IS frames over tcp from the specified server.
// Filter(s) are optional and use the following syntax:
//
// http://www.aprs-is.net/javAPRSFilter.aspx
func RecvIS(ctx context.Context, dial string, user Addr, pass int, filters ...string) <-chan Frame {
	fc := make(chan Frame)

	go func() {
		defer close(fc)

		conn, err := net.Dial("tcp", dial)
		if err != nil {
			return
		}
		defer conn.Close()

		// Read welcome banner
		_, err = readLine(conn, 5*time.Second)
		if err != nil {
			return
		}

		// Login
		login := genLogin(user, pass)
		if len(filters) > 0 {
			login += " filter " + strings.Join(filters, " ")
		}
		_, err = fmt.Fprintf(conn, "%s\r\n", login)
		if err != nil {
			return
		}
		// # logresp CWxxxx unverified, server CWOP-7
		// # logresp CWxxxx unverified, server THIRD
		_, err = readLine(conn, 5*time.Second)
		if err != nil {
			return
		}

		// Listen for frames until either the connection is closed or a
		// context cancel is received.
		var s string
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Heartbeats come across every 20 seconds so that's the
			// longest the read should block.  It's also the longest
			// it would take for a context cancel to be processed.
			s, err = readLine(conn, 30*time.Second)
			if err != nil {
				return
			}

			// # aprsc 2.1.4-g408ed49 26 Aug 2017 16:49:48 GMT FIFTH 44.74.128.25:14580
			if !strings.HasPrefix(s, "#") {
				f := Frame{}
				err = f.FromString(s)
				if err != nil {
					continue
				}
				fc <- f
			}
		}
	}()

	return fc
}

// SendIS sends a Frame to the specified APRS-IS dial string.  The
// dial string should be in the form scheme://host:port with
// scheme being http, tcp, or udp.  This is most commonly used for
// CWOP.
func (f Frame) SendIS(dial string, user Addr, pass int) error {
	// Refer to Connecting to APRS-IS:
	// http://www.aprs-is.net/Connecting.aspx

	parts := strings.Split(strings.ToLower(dial), "://")
	if len(parts) != 2 {
		return net.InvalidAddrError(dial)
	}

	switch parts[0] {
	case "http":
		return f.SendHTTP(dial, user, pass)
	case "tcp":
		return f.SendTCP(parts[1], user, pass)
	case "udp":
		return f.SendUDP(parts[1], user, pass)
	}

	return ErrProtoScheme
}

// SendHTTP sends a Frame to the specified APRS-IS host over the
// HTTP protocol.  This scheme is the least efficient and requires
// a verified connection (real callsign and passcode) but is
// reliable and provides acknowledgement of receipt.
func (f Frame) SendHTTP(dial string, user Addr, pass int) (err error) {
	if pass < 0 {
		err = ErrCallNotVerified
		return
	}

	data := fmt.Sprintf("%s\r\n%s", genLogin(user, pass), f)

	var req *http.Request
	req, err = http.NewRequest("POST", dial, bytes.NewBufferString(data))
	if err != nil {
		return
	}
	req.Header.Set("Accept-Type", "text/plain")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", strconv.Itoa(len(data)))

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request returned non-OK status code %d", resp.StatusCode)
	}

	return
}

// SendUDP sends a Frame to the specified APRS-IS host over the
// UDP protocol.  This scheme is the most efficient but requires
// a verified connection (real callsign and passcode) and has no
// acknowledgement of receipt.
func (f Frame) SendUDP(dial string, user Addr, pass int) (err error) {
	if pass < 0 {
		err = ErrCallNotVerified
		return
	}

	var conn net.Conn
	conn, err = net.Dial("udp", dial)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send data packet
	_, err = fmt.Fprintf(conn, "%s\r\n%s", genLogin(user, pass), f)

	return
}

// SendTCP sends a Frame to the specified APRS-IS host over the
// TCP protocol.  This scheme is the oldest, most compatible, and
// allows unverified connections.
func (f Frame) SendTCP(dial string, user Addr, pass int) (err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", dial)
	if err != nil {
		return
	}
	defer conn.Close()

	// Read welcome banner
	response, err := readLine(conn, 5*time.Second)
	if err != nil {
		return
	}
	log.Println("IN ", response)

	// Login
	_, err = fmt.Fprintf(conn, "%s\r\n", genLogin(user, pass))
	if err != nil {
		return
	}
	log.Printf("OUT %s\r\n", genLogin(user, pass))

	// # logresp CWxxxx unverified, server CWOP-7
	// # logresp CWxxxx unverified, server THIRD
	response, err = readLine(conn, 5*time.Second)
	if err != nil {
		return
	}
	log.Println("IN ", response)

	// Send frame
	_, err = fmt.Fprintf(conn, "%s\r\n", f)
	log.Printf("OUT %s\r\n", f)

	return
}

func processTcpSocket(conn net.Conn, frames chan Frame) {
	ticker := time.NewTicker(1 * time.Second)

readWriteLoop:
	for {
		select {
		// Check if there are messages to send
		case f := <-frames:
			// Send frame
			_, err := fmt.Fprintf(conn, "%s\r\n", f)
			log.Printf("OUT %s\r\n", f)
			if err != nil {
				log.Println("Socket error (write)", err.Error())
				break readWriteLoop
			}

		// Periodically check for data to read
		case _ = <-ticker.C:
			s, err := readLine(conn, 10*time.Millisecond)
			if err != nil {
				if err == io.EOF {
					log.Println("Socket closed (read)", err.Error())
					break readWriteLoop
				} else if err, ok := err.(net.Error); ok && err.Timeout() {
					// Under normal circumstance we will get a timeout if nothing was received
					//log.Println(".")
				} else {
					log.Println("Socket error (read)", err.Error())
					break readWriteLoop
				}
			}

			if len(s) > 0 {
				log.Println("IN ", s)
			}
		}
	}
}

// SendTCPFromChannel sends Frames to the specified APRS-IS host over the
// TCP protocol. The TCP socket is kept open, and reconnected on failure.
// Frames to send are provided via a channel.
func SendTCPFromChannel(dial string, user Addr, pass int, frames chan Frame) {
	for {
		conn, err := net.Dial("tcp", dial)
		fmt.Print("connect (", dial)
		if err != nil {
			fmt.Println(") fail")
		} else {
			fmt.Println(") ok")
			defer conn.Close()

			// Read welcome banner
			response, err := readLine(conn, 5*time.Second)
			if err != nil {
				return
			}
			log.Println("IN ", response)

			// Login
			_, err = fmt.Fprintf(conn, "%s\r\n", genLogin(user, pass))
			if err != nil {
				return
			}
			log.Printf("OUT %s\r\n", genLogin(user, pass))

			// # logresp CWxxxx unverified, server CWOP-7
			// # logresp CWxxxx unverified, server THIRD
			response, err = readLine(conn, 5*time.Second)
			if err != nil {
				return
			}
			log.Println("IN ", response)

			processTcpSocket(conn, frames)
		}
		time.Sleep(60 * time.Second) // TODO use https://pkg.go.dev/github.com/cenkalti/backoff/v4
	}
}
