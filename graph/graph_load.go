package main

import (
	"log"
	"time"

	"github.com/caglar10ur/rrd"
)

func main() {
	const dbfile = "/home/caglar/raspik/raspik.rrd"

	g := rrd.NewGrapher()

	g.SetVLabel("System Load")
	g.SetSize(800, 300)

	g.SetWatermark(time.Now().Format(time.RFC822))
	g.SetAltAutoscaleMax()
	g.SetSlopeMode()

	g.SetInterlaced()
	g.SetImageFormat("PNG")

	g.Def("1min", dbfile, "One", "AVERAGE")
	g.Def("5min", dbfile, "Five", "AVERAGE")
	g.Def("15min", dbfile, "Fifteen", "AVERAGE")

	// get the average
	g.VDef("1minAverage", "1min,AVERAGE")
	g.VDef("5minAverage", "5min,AVERAGE")
	g.VDef("15minAverage", "15min,AVERAGE")

	// use the average if the value is UN
	g.CDef("1minModified", "1min,UN,1minAverage,1min,IF")
	g.CDef("5minModified", "5min,UN,5minAverage,5min,IF")
	g.CDef("15minModified", "15min,UN,15minAverage,15min,IF")

	g.Line(1, "1minModified", "33cc33", "1 Min Load Avg")
	g.Line(1, "5minModified", "ff0000", "5 Min Load Avg")
	g.Line(1, "15minModified", "0000ff", "15 Min Load Avg")

	now := time.Now()

	g.SetTitle("Last 24 Hours")
	_, err := g.SaveGraph("/home/caglar/raspik/systemload24h.png", now.Add(-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Week")
	_, err = g.SaveGraph("/home/caglar/raspik/systemload7d.png", now.Add(7*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Month")
	_, err = g.SaveGraph("/home/caglar/raspik/systemload30d.png", now.Add(30*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

}
