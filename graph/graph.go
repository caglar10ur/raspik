package main

import (
	"log"
	"time"

	"github.com/caglar10ur/rrd"
)

func main() {
	const dbfile = "load.rrd"

	g := rrd.NewGrapher()

	g.SetTitle("Last 24 Hours")
	g.SetVLabel("System Load")
	g.SetSize(800, 300)

	g.SetWatermark(time.Now().Format(time.RFC822))
	g.SetAltAutoscaleMax()
	g.SetSlopeMode()

	g.SetInterlaced()
	g.SetImageFormat("PNG")

	g.Def("1min", dbfile, "load1", "AVERAGE")
	g.Def("5min", dbfile, "load5", "AVERAGE")
	g.Def("15min", dbfile, "load15", "AVERAGE")

	g.Line(1, "1min", "33cc33", "1 Min Load Avg")
	g.Line(1, "5min", "ff0000", "5 Min Load Avg")
	g.Line(1, "15min", "0000ff", "15 Min Load Avg")

	now := time.Now()
	_, err := g.SaveGraph("/home/caglar/web/rrd.png", now.Add(-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}
}
