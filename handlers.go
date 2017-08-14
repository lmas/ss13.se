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
	var hub ServerEntry
	if index > -1 {
		hub = servers[index]
		servers = append(servers[:index], servers[index+1:]...)
	}

	return a.templates["index"].Execute(w, map[string]interface{}{
		"Servers":      servers,
		"Hub":          hub,
		"TotalServers": len(servers),
	})
}

func (a *App) pageNews(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	return a.templates["news"].Execute(w, map[string]interface{}{
		"Reddit": a.news,
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

	c := makeHistoryChart("Daily history", points)
	return a.renderChart(w, c)
}

func (a *App) pageWeeklyChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
	id := vars["id"]
	points, err := a.store.GetSingleServerHistory(id, 7)
	if err != nil {
		return err
	}
	if len(points) < 1 {
		return HttpError{
			Status: 404,
			Err:    fmt.Errorf("server not found"),
		}
	}

	c := makeHistoryChart("Weekly history", points)
	return a.renderChart(w, c)
}

func (a *App) pageAverageChart(w http.ResponseWriter, r *http.Request, vars handlerVars) error {
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

	c := makeDayAverageChart("Average per day", points)
	return a.renderChart(w, c)
}
