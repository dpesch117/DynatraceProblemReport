package chart

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func BarChart(problemData map[string][]string) {
	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Problems sorted by Management Zones",
		Subtitle: "Default is set to past 2 hours",
	}))

	barItems := make([]opts.BarData, 0)
	var managementZones []string

	for key, element := range problemData {
		managementZones = append(managementZones, key)
		barItems = append(barItems, opts.BarData{Value: element[0]})
	}

	// Put data into instance
	bar.SetXAxis(managementZones).AddSeries("Category A", barItems)
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)
}
