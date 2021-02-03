package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/components"
	"github.com/go-echarts/go-echarts/opts"
)

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "Problems"},
		{Name: "Infrastructure"},
		{Name: "Services"},
		{Name: "Application"},
		{Name: "1Platform"},
		{Name: "Digital"},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "Problems", Target: "Infrastructure", Value: 10},
		{Source: "Problems", Target: "Services", Value: 18},
		{Source: "Problems", Target: "Application", Value: 20},
		{Source: "Infrastructure", Target: "1Platform", Value: 10},
		{Source: "Services", Target: "1Platform", Value: 15},
		{Source: "Services", Target: "Digital", Value: 3},
		{Source: "Application", Target: "Digital", Value: 20},
	}
)

func sankeyBase() *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Sankey-basic-example",
		}),
	)

	sankey.AddSeries("sankey", sankeyNode, sankeyLink, charts.WithLabelOpts(opts.Label{Show: true}))
	return sankey
}

type SankeyExamples struct{}

func (SankeyExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		sankeyBase(),
	)

	f, err := os.Create("examples/html/sankey.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
