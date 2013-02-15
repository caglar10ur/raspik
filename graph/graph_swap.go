package main

import (
	"log"
	"time"

	"github.com/caglar10ur/rrd"
)

func main() {
	const dbfile = "/home/caglar/raspik/raspik.rrd"

	g := rrd.NewGrapher()

	g.SetVLabel("System Swap")
	g.SetSize(800, 300)

	g.SetWatermark(time.Now().Format(time.RFC822))
	g.SetAltAutoscaleMax()
	g.SetSlopeMode()

	g.SetInterlaced()
	g.SetImageFormat("PNG")

	g.Def("total", dbfile, "TotalSwap", "AVERAGE")
	g.Def("used", dbfile, "UsedSwap", "AVERAGE")
	g.Def("free", dbfile, "FreeSwap", "AVERAGE")

	g.VDef("totalAverage", "total,AVERAGE")
	g.VDef("usedAverage", "used,AVERAGE")
	g.VDef("freeAverage", "free,AVERAGE")

	g.CDef("totalModified", "total,UN,totalAverage,total,IF")
	g.CDef("usedModified", "used,UN,usedAverage,used,IF")
	g.CDef("freeModified", "free,UN,freeAverage,free,IF")

	g.CDef("usedModifiedNeg", "0,usedModified,-")

	g.Line(2, "totalModified", "FF0000", "Total Swap")
	g.Area("usedModifiedNeg", "0000FF", "Used Swap")
	g.Area("freeModified", "00FF00", "Free Swap")

	g.HRule("0", "000000")

	now := time.Now()
	g.SetTitle("Last 24 Hours")
	_, err := g.SaveGraph("/home/caglar/raspik/systemswap24h.png", now.Add(-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Week")
	_, err = g.SaveGraph("/home/caglar/raspik/systemswap7d.png", now.Add(7*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}

	g.SetTitle("Last Month")
	_, err = g.SaveGraph("/home/caglar/raspik/systemswap30d.png", now.Add(30*-24*time.Hour), now)
	if err != nil {
		log.Fatal(err)
	}
}
