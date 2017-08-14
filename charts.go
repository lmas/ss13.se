package ss13_se

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

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

func makeHistoryChart(title string, points []ServerPoint) chart.Chart {
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
	li := &chart.LinearRegressionSeries{
		Name:        "Linear regression",
		InnerSeries: series,
	}
	sma := &chart.SMASeries{
		Name:        "Simple moving avg.",
		InnerSeries: series,
	}

	c := chart.Chart{
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Left: 120,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			series,
			li,
			sma,
		},
	}
	// Add a legend
	c.Elements = []chart.Renderable{
		chart.LegendLeft(&c),
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

	prettyName := func(d time.Weekday) string {
		return fmt.Sprintf("%s (%.0f)", d, avgDays[d])
	}

	return chart.BarChart{
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: []chart.Value{
			chart.Value{
				Label: prettyName(time.Monday),
				Value: avgDays[time.Monday],
			},
			chart.Value{
				Label: prettyName(time.Tuesday),
				Value: avgDays[time.Tuesday],
			},
			chart.Value{
				Label: prettyName(time.Wednesday),
				Value: avgDays[time.Wednesday],
			},
			chart.Value{
				Label: prettyName(time.Thursday),
				Value: avgDays[time.Thursday],
			},
			chart.Value{
				Label: prettyName(time.Friday),
				Value: avgDays[time.Friday],
			},
			chart.Value{
				Label: prettyName(time.Saturday),
				Value: avgDays[time.Saturday],
			},
			chart.Value{
				Label: prettyName(time.Sunday),
				Value: avgDays[time.Sunday],
			},
		},
	}
}
