package corvallisbus

import (
	"encoding/json"
	"fmt"
	cts "github.com/cvanderschuere/go-connexionz"
	"net/http"

	"appengine"
	"strconv"
)

const baseURL = "http://www.corvallistransit.com/"

func init() {
	http.HandleFunc("/platforms", platforms)
	http.HandleFunc("/routes", routes)
	http.HandleFunc("/patterns", patterns)
	http.HandleFunc("/cron/init", CTS_InfoInit)
}

func routes(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	c := cts.New(context, baseURL)

	//Create Platform (Need tag or number)
	intVal, _ := strconv.ParseInt(r.FormValue("tag"), 10, 64)
	p := &cts.Platform{
		Tag: intVal,
	}

	//Update ETA
	routes, err := c.ETA(p)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Return as json
	data, errJSON := json.Marshal(routes)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

func platforms(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	//Query for all platforms
	c := cts.New(context, baseURL)
	platforms, err := c.Platforms()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Return as json
	data, errJSON := json.Marshal(platforms)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

func patterns(w http.ResponseWriter, r *http.Request) {
	context := appengine.NewContext(r)

	c := cts.New(context, baseURL)

	routes, err := c.Patterns()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ps := make(map[string]*cts.Pattern)

	for _, route := range routes {
		if route.Number == r.FormValue("routeNum") || r.FormValue("routeNum") == "" {
			ps[route.Number] = route.Destination[0].Patterns[0]
		}
	}

	//Print out the pattern information
	data, errJSON := json.Marshal(ps)
	if errJSON != nil {
		http.Error(w, errJSON.Error(), 500)
		return
	}
	fmt.Fprint(w, string(data))
}
