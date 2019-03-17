package main

import (
	zmq "github.com/pebbe/zmq4"

	"github.com/caglar10ur/gologger"
	"github.com/caglar10ur/raspik"

	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"time"
)

type Stats struct {
	raspik.Load
	raspik.Uptime
	raspik.Mem
	raspik.Swap
}

var (
	Hostname string
	Port     uint
	Debug    bool
)

func init() {
	flag.UintVar(&Port, "Port", 5000, "Port number")
	flag.StringVar(&Hostname, "Hostname", "10ur.org", "Hostname")
	flag.BoolVar(&Debug, "Debug", false, "Debug")
	flag.Parse()
}
func main() {
	var stat Stats

	// Stand-in for a network connection
	var network bytes.Buffer

	log := logger.New(nil)
	if Debug {
		log.SetLogLevel(logger.Debug)
	}

	// types
	load := raspik.Load{}
	up := raspik.Uptime{}
	mem := raspik.Mem{}
	swap := raspik.Swap{}

	// polymorphism
	getter := [...]raspik.Getter{&load, &up, &mem, &swap}

	// socket
	socket, _ := zmq.NewSocket(zmq.PUB)
	defer socket.Close()

	socket.Connect(fmt.Sprintf("tcp://%s:%d", Hostname, Port))

	for {
		// collect values
		for _, g := range getter {
			g.Get()
		}

		// stat
		stat = Stats{Load: load, Uptime: up, Mem: mem, Swap: swap}

		// encode stat into network
		enc := gob.NewEncoder(&network)
		enc.Encode(stat)

		log.Debugf("Sending %+v\n", stat)

		// send it as multi-part msg topic + stat
		socket.SendMessage("raspik", &network)

		// clear the buffer
		network.Reset()

		// sleep 30 sec
		time.Sleep(30 * time.Second)
	}
}
