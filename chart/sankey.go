package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "Problems"},
		{Name: "1Platform"},
		{Name: "1Platform - DEV"},
		{Name: "Host"},
		{Name: "Service"},
		{Name: "Application"},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "Problems", Target: "1Platform", Value: 30},
		{Source: "Problems", Target: "1Platform - DEV", Value: 15},
		{Source: "1Platform", Target: "Host", Value: 10},
		{Source: "1Platform - DEV", Target: "Service", Value: 15},
		{Source: "1Platform", Target: "Service", Value: 20},
	}
)

func sankeyBase() *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Dynatrace Problem Distribution",
		}),
	)

	sankey.AddSeries("sankey", sankeyNode, sankeyLink, charts.WithLabelOpts(opts.Label{Show: true}))
	return sankey
}

func Sankey() {

	page := components.NewPage()
	page.AddCharts(
		sankeyBase(),
	)

	f, err := os.Create("sankey.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
