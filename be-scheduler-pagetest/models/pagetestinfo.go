package models


// Example struct
type Result struct {
	LookupUrl string `json:"lookupUrl"`
	TestStatusCode int    `json:"testStatusCode"`
	StatusText string  `json:"statusText, omitempty"`
	PageLoadTime int   `json:"pageLoadTime"`
	TestResults string `json:"testResults, omitempty"`
	PerformanceScore int `json:"performanceScore"`
	Availability float64   `json:"availability,omitempty"`
	AppId string      `json:"appId"`
	ScoreDetailsResult ScoreDetailResult `json:"overallscore"`
	Runs int	  `json:"runs,omitempty"`
}

type scoreDetailResult struct {

}

type LocationsResponse struct {
	StatusCode int `json:"statusCode"`
	Test  []TestAgentInfo `json:"testAgents"`
}

type TestAgentInfo struct {
	TestLabel    string `json:"label"`
	TestLocation string `json:"location"`
	TestBrowser  string `json:"browser"`
}


type UrlStat struct {
	URL string `json:"url"`
	PerformanceScore int `json:"performanceScore"`
	PageLoadTime int   `json:"pageLoadTime"`
	PageLoadHistory []LoadtimeHistory `json:"pageLoadTimes"`
	PageScoreHistory []ScoreHistory `json:"performanceScores"`
	AverageLoadTime int `json:"averageLoadTime"`
	Availability float64 `json:"availability"`
	DashboardId  string  `json:"dashboardId"`
}

type TestOverview struct {
	Count int `json:"count"`
  UrlStats []UrlStat `json:"urlStats"`
}