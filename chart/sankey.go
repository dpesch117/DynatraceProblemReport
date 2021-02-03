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
		{Name: "category1"},
		{Name: "category2"},
		{Name: "category3"},
		{Name: "category4"},
		{Name: "category5"},
		{Name: "category6"},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "category1", Target: "category2", Value: 10},
		{Source: "category2", Target: "category3", Value: 15},
		{Source: "category3", Target: "category4", Value: 20},
		{Source: "category5", Target: "category6", Value: 25},
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

func (SankeyExamples) Sankey() {
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
