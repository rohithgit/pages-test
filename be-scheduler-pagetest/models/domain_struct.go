package models

type DomainResult struct {
	DomainName    string `json:"name"`
	LoadTime int `json:"loadtime"`
	PageSize int `json:"pagesize"`
	Requests int `json:"results"`
	Percent float64 `json:"percent"`
	CDN      string  `json:"cdn,omitempty"`
	Connections int  `json:"connections,omitempty"`
}
