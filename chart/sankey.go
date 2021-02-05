package chart

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "Problems"},
		{Name: "Infrastructure"},
		{Name: "Service"},
	}

	sankeyLink = []opts.SankeyLink{}
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

func Sankey(problemData map[string][]string) {

	var managementZones []string

	for key, element := range problemData {
		managementZones = append(managementZones, key)
		fmt.Println("key: ", key, "value: ", element)
		sankeyNode = append(sankeyNode, opts.SankeyNode{Name: key})
		fmt.Println(sankeyNode)

		valTotal, _ := strconv.ParseFloat(element[0], 32)
		floatTotalProblems := float32(valTotal)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: "Problems", Target: key, Value: floatTotalProblems})
		valInfra, _ := strconv.ParseFloat(element[1], 32)
		floatInfraProblems := float32(valInfra)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: key, Target: "Infrastructure", Value: floatInfraProblems})
		valService, _ := strconv.ParseFloat(element[2], 32)
		floatServiceProblems := float32(valService)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: key, Target: "Service", Value: floatServiceProblems})

		fmt.Println("printing sankeyLink", sankeyLink)
		fmt.Println("total", floatTotalProblems, "infra", floatInfraProblems, "service", floatServiceProblems)

	}

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
