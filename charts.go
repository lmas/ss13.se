package ss13_se

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

var weekDaysOrder = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

type renderableChart interface {
	Render(chart.RendererProvider, io.Writer) error
}

func (a *App) renderChart(w http.ResponseWriter, c renderableChart) error {
	buf := &bytes.Buffer{}
	err := c.Render(chart.PNG, buf)

	if err != nil {
		//a.Log("Error while rendering chart: %s", err)
		return HttpError{
			Status: http.StatusInternalServerError,
			Err:    fmt.Errorf("error while rendering chart"),
		}
	}

	w.Header().Add("Content-Type", "image/png")
	_, err = io.Copy(w, buf)
	if err != nil {
		a.Log("Error while sending chart: %s", err)
		return HttpError{
			Status: http.StatusInternalServerError,
			Err:    fmt.Errorf("error while sending chart"),
		}
	}

	return nil
}

func makeHistoryChart(points []ServerPoint, showLegend bool) chart.Chart {
	// TODO BUG: one day is missing randomly (usually the 3rd day in the range) in the chart
	var xVals []time.Time
	var yVals []float64
	for _, p := range points {
		xVals = append(xVals, p.Time)
		yVals = append(yVals, float64(p.Players))
	}

	series := chart.TimeSeries{
		Name:    "Players",
		XValues: xVals,
		YValues: yVals,
	}
	lr := &chart.LinearRegressionSeries{
		Name:        "Linear regression",
		InnerSeries: series,
	}
	sma := &chart.SMASeries{
		Name:        "Simple moving avg.",
		InnerSeries: series,
	}

	c := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				t := int64(v.(float64))
				return time.Unix(0, t).Format("Jan 02 15:04")
			},
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f", v)
			},
		},
		Series: []chart.Series{
			series,
			lr,
			sma,
		},
	}
	if showLegend {
		c.Elements = []chart.Renderable{
			chart.LegendThin(&c),
		}
	}
	return c
}

// NOTE: The chart won't be renderable unless we've got at least two days/hours of history
func makeAverageChart(values map[int][]int, fnFormat func(int, float64) string, fnSort func([]int) []int) chart.BarChart {
	var keys []int
	avg := make(map[int]float64)
	for i, vl := range values {
		sum := 0
		for _, v := range vl {
			sum += v
		}
		avg[i] = float64(sum / len(vl))
		keys = append(keys, i)
	}

	var bars []chart.Value
	for _, k := range fnSort(keys) {
		bars = append(bars, chart.Value{
			Label: fnFormat(k, avg[k]),
			Value: avg[k],
			Style: chart.Style{
				StrokeColor: chart.ColorBlue,
				FillColor:   chart.ColorBlue,
			},
		})
	}

	barW, barS := 50, 100
	if len(avg) > 7 {
		barW, barS = 20, 20
	}
	s := chart.Style{
		Show:        true,
		StrokeWidth: 1,
	}
	return chart.BarChart{
		BarWidth:   barW,
		BarSpacing: barS,
		XAxis:      s,
		YAxis: chart.YAxis{
			Style: s,
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f", v)
			},
		},
		Bars: bars,
	}
}

// Shortcut/helper func for the calling handler
func avgDailyChart(points []ServerPoint) chart.BarChart {
	days := make(map[int][]int)
	for _, p := range points {
		d := int(p.Time.Weekday())
		days[d] = append(days[d], p.Players)
	}
	now := time.Now()
	formatter := func(i int, f float64) string {
		d := time.Weekday(i)
		extra := ""
		if d == now.Weekday() {
			extra = "*"
		}
		return fmt.Sprintf("%s%s", d, extra)
	}
	sorter := func(keys []int) []int {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		// Fucking wankers and their fucking sundays
		if keys[0] == int(time.Sunday) {
			keys = append(keys[1:], int(time.Sunday))
		}
		return keys
	}
	return makeAverageChart(days, formatter, sorter)
}

// Shortcut/helper func for the calling handler
func avgHourlyChart(points []ServerPoint) chart.BarChart {
	hours := make(map[int][]int)
	for _, p := range points {
		h := p.Time.Hour()
		hours[h] = append(hours[h], p.Players)
	}
	now := time.Now()
	formatter := func(i int, f float64) string {
		extra := ""
		if i == now.Hour() {
			extra = "*"
		}
		return fmt.Sprintf("%02d%s", i, extra)
	}
	sorter := func(keys []int) []int {
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		return keys
	}
	return makeAverageChart(hours, formatter, sorter)
}
