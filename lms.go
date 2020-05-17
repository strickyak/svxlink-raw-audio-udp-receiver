// build +main

// lms.go: A calibration program for Rx volume in svxlink.
//
// Print 28 plus the Log-Mean-Square of groups of 50 samples,
// and then ":" and the max of every group of 10 of those.
// Listens on UDP for the packet format that comes out of
// svxlink when you specify RAW_AUDIO_UDP_DEST=127.0.0.1:1825
// in your [Rx1] clause.  The samples are between -1 and 1,
// so their squares are between 0 and 1.  The printed 28+LMS with my
// IC-2730 radio with volume at 1:30 and "Mic" device at 42% volume
// tend to be
//   -- 7 to 8 when radio squelch is closed,
//   -- 15 to 18 when radio squelch opens but has dead air (receiving empty carrier),
//   -- 19 to 23 when I am speaking,
//   -- 20 to 24.5 for the repeater voice ID,
//   -- 22.7 to 23.2 for the bee-boop,
//   -- 25 is just about the maximum,
//   -- 25.77 on the carrier drop.
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"time"
)

var LISTEN = flag.String("listen", "127.0.0.1:1825", "UDP port to listen on")

func main() {
	flag.Parse()

	listenOn, err := net.ResolveUDPAddr("udp", *LISTEN)
	if err != nil {
		log.Fatalf("cannot net.ResolveUDPAddr %q: %v", *LISTEN, err)
	}

	conn, err := net.ListenUDP("udp", listenOn)
	if err != nil {
		log.Fatalf("cannot net.ListenUDP %q: %v", *LISTEN, err)
	}
	defer conn.Close()

	bb := make([]byte, 2000)
	packet := 0
	for {
		size, fromWhere, err := conn.ReadFromUDP(bb)
		if err != nil {
			log.Printf("cannot conn.ReadFromUDP %q: %v", *LISTEN, err)
		}
		_ = fromWhere
		now := float64(time.Now().UnixNano()) / 1000000000
		_ = now

		// We have been getting 2000 bytes, or 500 sample.
		// We see 47 of those per second.
		//    47 * 500 = 23500 hz (!?)
		// log.Printf("From %d got %d bytes, %d samples", fromWhere.Port, size, size/4);

		r := bytes.NewReader(bb[:size])
		var maxx float64
		for i := 0; i < size/4; i += 50 {
			var v float32
			var sumsq float64
			for j := i; j < i+50; j++ {
				err := binary.Read(r, binary.LittleEndian, &v)
				if err != nil {
					log.Fatalf("Cannot binary.Read: %v", err)
				}
				// fmt.Printf("%f ", v)
				sumsq += float64(v * v)
			}
			x := 28 + math.Log(sumsq/50)
			fmt.Printf("%6.2f ", x)
			if maxx < x {
				maxx = x
			}
		}
		fmt.Printf(": %6.2f (%d) %.3f\n", maxx, packet, now)
		// fmt.Println()
		// fmt.Printf("[%5d]%20.3f  n=%8g   %10.4f\n", packet, now, n, math.Log(sumsq/n))
		packet++
	}
}
