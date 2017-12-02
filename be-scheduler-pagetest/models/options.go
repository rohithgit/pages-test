package models

import (
	"encoding/json"
	"fmt"
	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/spec-ap/scheduler-pagetest/data"
)

// Options contains the global options for the template service.
type Options struct {
	Protocol         string     `json:"protocol,omitempty"`
	Debug            bool       `json:"debug,omitempty"`
	Verbose          bool       `json:"verbose,omitempty"`
	Port             int        `json:"port,omitempty"`
	Mock             bool       `json:"mock,omitempty"`
	Environment      string     `json:"environment"`
	CertFile         string     `json:"certFile,omitempty"`
	PkeyFile         string     `json:"pkeyFile,omitempty"`
	RedisServers     string     `json:"redis_servers"`
	WptServer 	 string     `json:"wpt_server"`
	Mongo 		 string     `json:"mongo"`
	IsMongoSecure    bool       `json:"ismongosecure,omitempty"`
	EnableSvcDiscovery bool   `json:"servicediscovery"`
	SchedulerEnabled 		 bool     `json:"enable,omitempty"`
	SchedulerInt 		 int     `json:"schedulerint"`
	SchedulerSleep           int     `json:"schedulersleep"`
	SecretKey string `json:"-"`
	ValidateToken bool   `json:"validatetoken,omitempty"`
	ClientId      string `json:"clientid,omitempty"`
	ClientSecret  string `json:"clientsecret,omitempty"`
	Elastic       string `json:"elastic,omitempty"`
	EsUser       string `json:"esuser,omitempty"`
	EsPassword       string `json:"espassword,omitempty"`
}

// NewOptions returns a ptr to a new options object
// the options are initialized from the default options
// json object in the data package
func NewOptions() *Options {
	return &Options{}
}

// InitOptions initializes the options
func InitOptions() (*Options, error) {
	// init service options
	options := NewOptions()
	if err := json.Unmarshal(data.DefaultOptions, options); err != nil {
		return nil, fmt.Errorf("Options initialization unmarshal error: %v", err)
	}
	return options, nil
}
