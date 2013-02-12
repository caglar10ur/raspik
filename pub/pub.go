package main

import (
	zmq "github.com/alecthomas/gozmq"

	"github.com/caglar10ur/gologger"
	"github.com/caglar10ur/raspik"

	"bytes"
	"encoding/gob"
	"time"
)

type Stats struct {
	raspik.Load
	raspik.Uptime
	raspik.Mem
	raspik.Swap
}

func main() {
	var stat Stats

	// Stand-in for a network connection
	var network bytes.Buffer

	log := logger.New(nil)
	log.SetLogLevel(logger.Debug)

	// types
	load := raspik.Load{}
	up := raspik.Uptime{}
	mem := raspik.Mem{}
	swap := raspik.Swap{}

	// polymorphism
	getter := [...]raspik.Getter{&load, &up, &mem, &swap}

	// context
	context, _ := zmq.NewContext()
	defer context.Close()

	// socket
	socket, _ := context.NewSocket(zmq.PUB)
	socket.Connect("tcp://10ur.org:5000")
	socket.SetSockOptUInt64(zmq.HWM, 1)
	defer socket.Close()

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
		socket.Send([]byte("raspik"), zmq.SNDMORE)
		socket.Send(network.Bytes(), 0)

		// clear the buffer
		network.Reset()

		// sleep 30 sec
		time.Sleep(30 * time.Second)
	}
}
