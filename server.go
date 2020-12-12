package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lkslts64/inaccessproject/ptask"
)

const layout = "20060102T150405Z"

var addr = flag.String("addr", ":8080", "address to listen (provide it like the net package specifies)")

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welocme to periodic tasks website")
}

func tasksHandler(w http.ResponseWriter, r *http.Request, period string) {
	q := r.URL.Query()
	enc := json.NewEncoder(w)
	tz, st1, st2 := q.Get("tz"), q.Get("t1"), q.Get("t2")
	t1, err1 := parseTime(layout, st1, tz)
	t2, err2 := parseTime(layout, st2, tz)
	if err2 != nil || err1 != nil {
		err := enc.Encode(&reqErr{
			Status: "error",
			Desc: fmt.Sprintf("malformed query values: (%s)", func(e1, e2 error) string {
				if e1 != nil {
					return e1.Error()
				}
				return e2.Error()
			}(err1, err2)),
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
	ptlist, err := ptask.List(t1, t2, period, 1<<20)
	if err != nil {
		err := enc.Encode(reqErr{
			Status: "error",
			Desc:   err.Error(),
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
	fptlist := make([]string, len(ptlist))
	for i := range ptlist {
		fptlist[i] = ptlist[i].UTC().Format(layout)
	}
	enc.Encode(fptlist)
}

func hourHandler(w http.ResponseWriter, r *http.Request) {
	tasksHandler(w, r, "1h")
}

func dayHandler(w http.ResponseWriter, r *http.Request) {
	tasksHandler(w, r, "1d")
}

func monthHandler(w http.ResponseWriter, r *http.Request) {
	tasksHandler(w, r, "1mo")
}

func yearHandler(w http.ResponseWriter, r *http.Request) {
	tasksHandler(w, r, "1y")
}

// parses `value` based on `layout` and returns a time.Time whose's location is
// set as `name`.
func parseTime(layout, value, name string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

type reqErr struct {
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

func (re *reqErr) Error() string {
	return fmt.Sprintf("status: %s, description: %s", re.Status, re.Desc)
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/", homePageHandler)
	r.HandleFunc("/tasks/1h", hourHandler).Methods("GET")
	r.HandleFunc("/tasks/1d", dayHandler).Methods("GET")
	r.HandleFunc("/tasks/1mo", monthHandler).Methods("GET")
	r.HandleFunc("/tasks/1y", yearHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(*addr, r))
}
