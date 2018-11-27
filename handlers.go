package ss13_se

import (
	"fmt"
	"net/http"
)

func (a *App) pageIndex(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	servers, err := a.store.GetServers()
	if err != nil {
		return err
	}

	// Remove the internal entry used to count total players
	index := -1
	for i, s := range servers {
		if s.Title == internalServerTitle {
			index = i
			break
		}
	}
	if index > -1 {
		servers = append(servers[:index], servers[index+1:]...)
	}

	return a.templates["index"].Execute(w, map[string]interface{}{
		"Servers": servers,
		"Hub":     a.hub,
	})
}

func (a *App) pageStyle(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	w.Header().Set("Content-Type", "text/css")
	return a.templates["style"].Execute(w, nil)
}

func (a *App) pageNews(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	return a.templates["news"].Execute(w, map[string]interface{}{
		"Reddit": a.news,
		"Hub":    a.hub,
	})
}

func (a *App) pageServer(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	server, err := a.store.GetServer(id)
	if err != nil {
		// TODO: handle and log the error properly
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	if server.Title == internalServerTitle {
		server.Title = "Global stats"
	}

	return a.templates["server"].Execute(w, map[string]interface{}{
		"Server": server,
		"Hub":    a.hub,
	})
}

func (a *App) pageDailyChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	points, err := a.store.GetSingleServerHistory(id, 1)
	if err != nil {
		return err
	}
	if len(points) < 1 {
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	c := makeHistoryChart(points, true)
	return a.renderChart(w, c)
}

func (a *App) pageWeeklyChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	points, err := a.store.GetSingleServerHistory(id, 6)
	if err != nil {
		return err
	}
	if len(points) < 1 {
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	c := makeHistoryChart(points, false)
	return a.renderChart(w, c)
}

func (a *App) pageAverageDailyChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	points, err := a.store.GetSingleServerHistory(id, 30)
	if err != nil {
		return err
	}
	if len(points) < 1 {
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	c := avgDailyChart(points)
	return a.renderChart(w, c)
}

func (a *App) pageAverageHourlyChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	points, err := a.store.GetSingleServerHistory(id, 30)
	if err != nil {
		return err
	}
	if len(points) < 1 {
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	c := avgHourlyChart(points)
	return a.renderChart(w, c)
}
