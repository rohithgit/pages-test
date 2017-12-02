// Package router routes http REST calls
package router

import (
	"github.com/gorilla/mux"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/handlers"
)

// NewTopoRouter returns a new router for the topo micro service
func NewTopoRouter() (*mux.Router, error) {
	r := mux.NewRouter()
	s := r.PathPrefix("/spectre/v1/pagetest").Methods("GET", "POST").Subrouter()

	// subrouter
	s.HandleFunc("/test", handlers.PageTest)
	s.HandleFunc("/result", handlers.LookupResults)

	s.HandleFunc("/pageloadtime/domains", handlers.DomainLoadTime)
	s.HandleFunc("/result/pageloadtime/domains", handlers.DomainLoadTimeResult)

	s.HandleFunc("/pageloadtime/contenttypes", handlers.ContentTypes)
	s.HandleFunc("/result/pageloadtime/contenttypes", handlers.ContentTypesResult)

	s.HandleFunc("/pageloadtime/states", handlers.States)
	s.HandleFunc("/result/pageloadtime/states", handlers.StatesResult)

	s.HandleFunc("/pageloadtime", handlers.PageLoad)
	s.HandleFunc("/result/pageloadtime", handlers.LookupPageLoad)

	s.HandleFunc("/pageloadtime/history", handlers.PageLoadChart)

	s.HandleFunc("/perfscore", handlers.PerfScore)
	s.HandleFunc("/result/perfscore", handlers.LookupPerfScore)

	s.HandleFunc("/perfscore/details", handlers.PerfScoreDetails)
	s.HandleFunc("/result/perfscore/details", handlers.LookupPerfScoreDetails)

	s.HandleFunc("/successfulruns", handlers.SuccessfulRuns)
	s.HandleFunc("/slowestrun", handlers.SlowestRun)

	s.HandleFunc("/locations", handlers.Locations)
	s.HandleFunc("/performanceavailability", handlers.Availability)
	s.HandleFunc("/hello", handlers.Hello)
	s.HandleFunc("/webendpoints", handlers.GetPageTestsForApp)
	s.HandleFunc("/webendpoints/count", handlers.GetPageTestsForAppCount)
	s.HandleFunc("/webendpoints/exists", handlers.URLExists)

	return r, nil
}
