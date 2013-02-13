package main

import (
	"log"
	"time"

	"github.com/caglar10ur/rrd"
)

func main() {
	const dbfile = "/home/caglar/raspik/raspik.rrd"

	g := rrd.NewGrapher()

	g.SetVLabel("System Memory")
	g.SetSize(800, 300)

	g.SetWatermark(time.Now().Format(time.RFC822))
	g.SetAltAutoscaleMax()
	g.SetSlopeMode()

	g.SetInterlaced()
	g.SetImageFormat("PNG")

	g.Def("total", dbfile, "TotalRam", "AVERAGE")
	g.Def("free", dbfile, "FreeRam", "AVERAGE")
	g.Def("shared", dbfile, "SharedRam", "AVERAGE")
	g.Def("buffer", dbfile, "BufferRam", "AVERAGE")

	g.VDef("totalAverage", "total,AVERAGE")
	g.VDef("freeAverage", "free,AVERAGE")
	g.VDef("sharedAverage", "shared,AVERAGE")
	g.VDef("bufferAverage", "buffer,AVERAGE")

	g.CDef("totalModified", "total,UN,totalAverage,total,IF")
	g.CDef("freeModified", "free,UN,freeAverage,free,IF")
	g.CDef("sharedModified", "shared,UN,sharedAverage,shared,IF")
	g.CDef("bufferModified", "buffer,UN,bufferAverage,buffer,IF")

	g.Line(1, "totalModified", "33cc33", "Total Ram")
	g.Line(1, "freeModified", "ff0000", "Free Ram")
	g.Line(1, "sharedModified", "0000ff", "Shared Ram")
	g.Line(1, "bufferModified", "996633", "Buffer")

	now := time.Now()
	g.SetTitle("Last 24 Hours")
	_, err := g.SaveGraph("/home/caglar/raspik/systemmem24h.png", now.Add(-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Week")
	_, err = g.SaveGraph("/home/caglar/raspik/systemmem7d.png", now.Add(7*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Month")
	_, err = g.SaveGraph("/home/caglar/raspik/systemmem30d.png", now.Add(30*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}
}
