package ss13_se

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func makeHistoryChart(title string, showLegend bool, points []ServerPoint) chart.Chart {
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
		// Add a legend
		c.Elements = []chart.Renderable{
			chart.LegendThin(&c),
		}
	}
	return c
}

// NOTE: The chart won't be renderable unless we've got at least two days of history
func makeDayAverageChart(title string, points []ServerPoint) chart.BarChart {
	days := make(map[time.Weekday][]int)
	for _, p := range points {
		day := p.Time.Weekday()
		days[day] = append(days[day], p.Players)
	}

	avgDays := make(map[time.Weekday]float64)
	for day, vals := range days {
		sum := 0
		for _, v := range vals {
			sum += v
		}
		avg := sum / len(vals)
		avgDays[day] = float64(avg)
	}

	var bars []chart.Value
	for _, d := range weekDaysOrder {
		bars = append(bars, chart.Value{
			Label: fmt.Sprintf("%s (%.0f)", d, avgDays[d]),
			Value: avgDays[d],
			Style: chart.Style{
				StrokeColor: chart.ColorBlue,
				FillColor:   chart.ColorBlue,
			},
		})
	}

	return chart.BarChart{
		BarWidth: 50,
		XAxis:    chart.StyleShow(),
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f", v)
			},
		},
		Bars: bars,
	}
}
