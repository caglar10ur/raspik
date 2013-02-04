package main

import (
	zmq "github.com/alecthomas/gozmq"

	"github.com/caglar10ur/raspik"
	"github.com/caglar10ur/rrd"

	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Eastern  = time.FixedZone("Eastern", -5*3600)
	Central  = time.FixedZone("Central", -6*3600)
	Mountain = time.FixedZone("Mountain", -7*3600)
	Pacific  = time.FixedZone("Pacific", -8*3600)
)

const (
	dbfile    = "load.rrd"
	step      = 30
	heartbeat = 2 * step
)

type Stats struct {
	raspik.Load
	raspik.Uptime
	raspik.Mem
	raspik.Swap
}

func createLoadRRD() error {
	c := rrd.NewCreator(dbfile, time.Now(), step)

	c.DS("load1", "GAUGE", heartbeat, 0, 100)
	c.DS("load5", "GAUGE", heartbeat, 0, 100)
	c.DS("load15", "GAUGE", heartbeat, 0, 100)

	// three RRAs with a resolution of 5 minutes spanning 31 days using the AVERAGE, MIN, and MAX consolidation functions,
	c.RRA("AVERAGE", 0.5, 1, 89280)
	c.RRA("MIN", 0.5, 1, 89280)
	c.RRA("MAX", 0.5, 1, 89280)
	// RRA with a resolution of 15 minutes for 90 days
	c.RRA("AVERAGE", 0.5, 3, 86400)
	// RRA with a resolution of 1 hour for 365 days.
	c.RRA("AVERAGE", 0.5, 12, 87600)

	// do not overwrite
	return c.Create(false)
}

func main() {
	var stat Stats

	// Stand-in for a network connection
	var network bytes.Buffer

	// context
	context, _ := zmq.NewContext()
	defer context.Close()

	// socket
	socket, _ := context.NewSocket(zmq.SUB)
	socket.Bind("tcp://*:5000")
	// filter out topics other than raspik
	socket.SetSockOptString(zmq.SUBSCRIBE, "raspik")
	defer socket.Close()

	// create RRD file if not exists
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		err = createLoadRRD()
		if err != nil {
			log.Fatal(err)
		}
	}

	// updater
	u := rrd.NewUpdater(dbfile)
	for {
		// receive multi-part msg topic + stat
		msg, err := socket.RecvMultipart(0)
		if err != nil {
			fmt.Printf("ERROR: %+v\n", msg)
		} else {
			// write to buffer
			network.Write(msg[1])

			// decode network into stat
			dec := gob.NewDecoder(&network)
			dec.Decode(&stat)

			// clear the buffer
			network.Reset()

			fmt.Printf("%s @ %+v\n", time.Now().In(Eastern).Format(time.RFC822), stat)

			// update RRD file
			err = u.Update(time.Now(), stat.Load.One, stat.Load.Five, stat.Load.Fifteen)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
