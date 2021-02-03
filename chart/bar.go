package chart

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func BarChart(managementZoneNames []string, problemCount []opts.BarData) {
	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Problems sorted by Management Zones",
		Subtitle: "Default is set to past 2 hours",
	}))

	// Put data into instance
	bar.SetXAxis(managementZoneNames).
		AddSeries("Category A", problemCount)
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)
}
