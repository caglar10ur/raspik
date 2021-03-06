package main

import (
	zmq "github.com/pebbe/zmq4"

	"github.com/caglar10ur/raspik"
	"github.com/ziutek/rrd"

	"github.com/caglar10ur/gologger"

	"bytes"
	"encoding/gob"
	"flag"
	"os"
	"time"
)

var (
	New_York, _ = time.LoadLocation("America/New_York")
	Debug       bool
)

func init() {
	flag.BoolVar(&Debug, "Debug", false, "Debug")
	flag.Parse()
}

const (
	dbfile    = "/home/caglar/raspik/raspik.rrd"
	step      = 30
	heartbeat = 2 * step
)

type Stats struct {
	raspik.Load
	raspik.Uptime
	raspik.Mem
	raspik.Swap
}

func createRRD() error {
	c := rrd.NewCreator(dbfile, time.Now(), step)

	// load
	c.DS("One", "GAUGE", heartbeat, 0, 100)
	c.DS("Five", "GAUGE", heartbeat, 0, 100)
	c.DS("Fifteen", "GAUGE", heartbeat, 0, 100)

	// mem
	c.DS("TotalRam", "GAUGE", heartbeat, 0, "U")
	c.DS("FreeRam", "GAUGE", heartbeat, 0, "U")
	c.DS("SharedRam", "GAUGE", heartbeat, 0, "U")
	c.DS("BufferRam", "GAUGE", heartbeat, 0, "U")

	// swap
	c.DS("TotalSwap", "GAUGE", heartbeat, 0, "U")
	c.DS("UsedSwap", "GAUGE", heartbeat, 0, "U")
	c.DS("FreeSwap", "GAUGE", heartbeat, 0, "U")

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

	log := logger.New(nil)
	if Debug {
		log.SetLogLevel(logger.Debug)
	}

	// socket
	socket, _ := zmq.NewSocket(zmq.SUB)
	socket.Bind("tcp://*:5000")
	// filter out topics other than raspik
	socket.SetSubscribe("raspik")
	defer socket.Close()

	// create RRD file if not exists
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		err = createRRD()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// updater
	u := rrd.NewUpdater(dbfile)
	for {
		// receive multi-part msg topic + stat
		msg, err := socket.RecvMessageBytes(0)
		if err != nil {
			log.Errorf("ERROR: %+v\n", msg)
		} else {
			// write to buffer
			network.Write(msg[1])

			// decode network into stat
			dec := gob.NewDecoder(&network)
			dec.Decode(&stat)

			// clear the buffer
			network.Reset()

			log.Debugf("%s @ %+v\n", time.Now().In(New_York).Format(time.RFC822), stat)

			// update RRD file
			err = u.Update(time.Now(), stat.Load.One, stat.Load.Five, stat.Load.Fifteen,
				stat.Mem.TotalRam, stat.Mem.FreeRam, stat.Mem.SharedRam, stat.Mem.BufferRam,
				stat.Swap.TotalSwap, stat.Swap.UsedSwap, stat.Swap.FreeSwap)
			if err != nil {
				log.Errorln(err)
			}
		}
	}
}
