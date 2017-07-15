package ss13_se

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HttpError struct {
	Status int
	Err    error
}

func (s HttpError) Error() string {
	return fmt.Sprintf("%d %s", s.Status, s.Err.Error())
}

type handlerVars map[string]string

type handler func(http.ResponseWriter, *http.Request, handlerVars) error

func (h handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	err := h(rw, req, mux.Vars(req))
	dur := time.Since(start)

	if err != nil {
		switch e := err.(type) {
		case HttpError:
			http.Error(rw, e.Error(), e.Status)
		default:
			http.Error(rw, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}

	log.Printf("%s %s \t%s \terr: %v\n",
		//req.RemoteAddr,
		req.Method,
		req.URL.String(),
		//req.UserAgent(),
		dur,
		//resp.Status,
		err,
	)
}
