package chart

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Sankeydata struct {
	Nodes []opts.SankeyNode `json:"nodes"`
	Links []opts.SankeyLink `json:"links"`
}

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

	//iterate over problemdata and append problems to sankeylink
	for key, element := range problemData {
		fmt.Println("key: ", key, "value: ", element)

		sankeyNode = append(sankeyNode, opts.SankeyNode{Name: key})

		valTotal, _ := strconv.ParseFloat(element[0], 32)
		floatTotalProblems := float32(valTotal)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: "Problems", Target: key, Value: floatTotalProblems})
	}

	//iterate over problemdata and append infrastructure to sankeylink
	for key, element := range problemData {
		valInfra, _ := strconv.ParseFloat(element[1], 32)
		floatInfraProblems := float32(valInfra)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: key, Target: "Infrastructure", Value: floatInfraProblems})
		//valService, _ := strconv.ParseFloat(element[2], 32)
		//floatServiceProblems := float32(valService)
		//sankeyLink = append(sankeyLink, opts.SankeyLink{Source: key, Target: "Service", Value: floatServiceProblems})
	}
	//iterate over problemdata and append Service to sankeylink
	for key, element := range problemData {
		valService, _ := strconv.ParseFloat(element[2], 32)
		floatServiceProblems := float32(valService)
		sankeyLink = append(sankeyLink, opts.SankeyLink{Source: key, Target: "Service", Value: floatServiceProblems})
	}

	//Creating JSON File
	sankeyData := Sankeydata{Nodes: sankeyNode, Links: sankeyLink}

	file, _ := json.MarshalIndent(sankeyData, "   ", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)

	fmt.Println("Printing Sankey Node")
	fmt.Println(sankeyNode)
	//fmt.Println(sankeyLink)

	page := components.NewPage()
	page.AddCharts(
		sankeyBase(),
		graphEnergy(),
	)

	f, err := os.Create("sankey.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func graphEnergy() *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Sankey-json-file-example",
		}),
	)

	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatal(err)
	}

	type Data struct {
		Nodes []opts.SankeyNode
		Links []opts.SankeyLink
	}

	var data Data
	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Println(err)
	}

	sankey.AddSeries("sankey", data.Nodes, data.Links).
		SetSeriesOptions(
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:     "source",
				Curveness: 0.5,
			}),
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
		)
	return sankey
}
