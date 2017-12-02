package models

import (
	"time"
	"encoding/json"
)

type SummaryData struct {
	TestId string `json:"testId"`
	OwnerKey string `json:"ownerKey"`
	JsonUrl string `json:"jsonUrl"`
	XmlUrl string `json:"xmlUrl"`
	UserUrl string `json:"userUrl"`
	SummaryCsv string `json:"summaryCSV"`
	DetailsCSV string `json:"detailsCSV"`
}

type Summary struct {
	StatusCode        int     `json:"statusCode"`
	StatusText       string     `json:"statusText"`
	Data       SummaryData     `json:"data"`
}

type DetailsData struct {
	Id string `json:"id"`
	Url string `json:"url"`
	Summary string `json:"summary"`
	TestUrl string `json:"testUrl"`
	Location string `json:"location"`
	From string `json:"from"`
	Connectivity string `json:"connectivity"`
	Median `json:"median"`
}

type Details struct {
	 DetailsData DetailsData `json:"data"`
}

type Median struct {
	FirstView FirstView `json:"firstView"`
}

type FirstView struct {
	LoadTime int `json:"loadTime"`
}


type JsonForTestStarted struct {
	Data struct {
		     Elapsed         int    `json:"elapsed"`
		     FvRunsCompleted int    `json:"fvRunsCompleted"`
		     Fvonly          int    `json:"fvonly"`
		     ID              string `json:"id"`
		     Location        string `json:"location"`
		     Remote          bool   `json:"remote"`
		     Runs            int    `json:"runs"`
		     RvRunsCompleted int    `json:"rvRunsCompleted"`
		     StartTime       string `json:"startTime"`
		     StatusCode      int    `json:"statusCode"`
		     StatusText      string `json:"statusText"`
		     TestID          string `json:"testId"`
		     TestInfo        struct {
					     Bodies       int    `json:"bodies"`
					     Browser      string `json:"browser"`
					     BwIn         int    `json:"bwIn"`
					     BwOut        int    `json:"bwOut"`
					     Connectivity string `json:"connectivity"`
					     Fvonly       int    `json:"fvonly"`
					     IgnoreSSL    int    `json:"ignoreSSL"`
					     Iq           int    `json:"iq"`
					     Keepua       int    `json:"keepua"`
					     Label        string `json:"label"`
					     Latency      int    `json:"latency"`
					     Location     string `json:"location"`
					     Mobile       int    `json:"mobile"`
					     Netlog       int    `json:"netlog"`
					     Noscript     int    `json:"noscript"`
					     Plr          string `json:"plr"`
					     Pngss        int    `json:"pngss"`
					     Priority     int    `json:"priority"`
					     Runs         int    `json:"runs"`
					     Scripted     int    `json:"scripted"`
					     Standards    int    `json:"standards"`
					     Tcpdump      int    `json:"tcpdump"`
					     Timeline     int    `json:"timeline"`
					     Trace        int    `json:"trace"`
					     URL          string `json:"url"`
					     Web10        int    `json:"web10"`
				     } `json:"testInfo"`
		     TestsCompleted int `json:"testsCompleted"`
		     TestsExpected  int `json:"testsExpected"`
		     URL	    string   `json:"url"`
	     } `json:"data"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

type JsonResults1 struct {
	Data struct {
		     Average struct {
				     FirstView struct {
						       SpeedIndex                                  int     `json:"SpeedIndex"`
						       TTFB                                        int     `json:"TTFB"`
						       AdultSite                                   int     `json:"adult_site"`
						       Aft                                         int     `json:"aft"`
						       AvgRun                                      int     `json:"avgRun"`
						       BytesIn                                     int     `json:"bytesIn"`
						       BytesInDoc                                  int     `json:"bytesInDoc"`
						       BytesOut                                    int     `json:"bytesOut"`
						       BytesOutDoc                                 int     `json:"bytesOutDoc"`
						       Cached                                      int     `json:"cached"`
						       ChromeUserTiming_domComplete                int     `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int     `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int     `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int     `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int     `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int     `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int     `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int     `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int     `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int     `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int     `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int     `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int     `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int     `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int     `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int     `json:"connections"`
						       Date                                        int     `json:"date"`
						       DocCPUms                                    float64 `json:"docCPUms"`
						       DocCPUpct                                   int     `json:"docCPUpct"`
						       DocTime                                     int     `json:"docTime"`
						       DomContentLoadedEventEnd                    int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart                  int     `json:"domContentLoadedEventStart"`
						       DomElements                                 int     `json:"domElements"`
						       DomTime                                     int     `json:"domTime"`
						       Domains                    map[string]struct {
							       Bytes       int `json:"bytes"`
							       Connections int `json:"connections"`
							       Requests    int `json:"requests"`
						       } `json:"domains"`
						       EffectiveBps                                int     `json:"effectiveBps"`
						       EffectiveBpsDoc                             int     `json:"effectiveBpsDoc"`
						       FirstPaint                                  int     `json:"firstPaint"`
						       FixedViewport                               int     `json:"fixed_viewport"`
						       FullyLoaded                                 int     `json:"fullyLoaded"`
						       FullyLoadedCPUms                            float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct                           int     `json:"fullyLoadedCPUpct"`
						       GzipSavings                                 int     `json:"gzip_savings"`
						       GzipTotal                                   int     `json:"gzip_total"`
						       ImageSavings                                int     `json:"image_savings"`
						       ImageTotal                                  int     `json:"image_total"`
						       IsResponsive                                int     `json:"isResponsive"`
						       LastVisualChange                            int     `json:"lastVisualChange"`
						       LoadEventEnd                                int     `json:"loadEventEnd"`
						       LoadEventStart                              int     `json:"loadEventStart"`
						       LoadTime                                    int     `json:"loadTime"`
						       MinifySavings                               int     `json:"minify_savings"`
						       MinifyTotal                                 int     `json:"minify_total"`
						       OptimizationChecked                         int     `json:"optimization_checked"`
						       PageSpeedVersion                            float64 `json:"pageSpeedVersion"`
						       Render                                      int     `json:"render"`
						       Requests                                    int     `json:"requests"`
						       RequestsDoc                                 int     `json:"requestsDoc"`
						       RequestsFull                                int     `json:"requestsFull"`
						       Responses200                                int     `json:"responses_200"`
						       Responses404                                int     `json:"responses_404"`
						       ResponsesOther                              int     `json:"responses_other"`
						       Result                                      int     `json:"result"`
						       Run                                         int     `json:"run"`
						       ScoreCache                                  int     `json:"score_cache"`
						       ScoreCdn                                    int     `json:"score_cdn"`
						       ScoreCombine                                int     `json:"score_combine"`
						       ScoreCompress                               int     `json:"score_compress"`
						       ScoreCookies                                int     `json:"score_cookies"`
						       ScoreEtags                                  int     `json:"score_etags"`
						       ScoreGzip                                   int     `json:"score_gzip"`
						       ScoreKeep_alive                             int     `json:"score_keep-alive"`
						       ScoreMinify                                 int     `json:"score_minify"`
						       ScoreProgressiveJpeg                        int     `json:"score_progressive_jpeg"`
						       ServerCount                                 int     `json:"server_count"`
						       ServerRtt                                   int     `json:"server_rtt"`
						       TitleTime                                   int     `json:"titleTime"`
						       VisualComplete                              int     `json:"visualComplete"`
					       } `json:"firstView"`
			     } `json:"average"`
		     BwDown       int    `json:"bwDown"`
		     BwUp         int    `json:"bwUp"`
		     Completed    int    `json:"completed"`
		     Connectivity string `json:"connectivity"`
		     From         string `json:"from"`
		     Fvonly       bool   `json:"fvonly"`
		     ID           string `json:"id"`
		     Latency      int    `json:"latency"`
		     Location     string `json:"location"`
		     Median       struct {
				     FirstView struct {
						       SpeedIndex  int    `json:"SpeedIndex"`
						       TTFB        int    `json:"TTFB"`
						       URL         string `json:"URL"`
						       AdultSite   int    `json:"adult_site"`
						       Aft         int    `json:"aft"`
						       BasePageCdn string `json:"base_page_cdn"`
						       Breakdown   interface{} `json:"breakdown"`
						       BrowserName      string `json:"browser_name"`
						       BrowserVersion   string `json:"browser_version"`
						       BytesIn          int    `json:"bytesIn"`
						       BytesInDoc       int    `json:"bytesInDoc"`
						       BytesOut         int    `json:"bytesOut"`
						       BytesOutDoc      int    `json:"bytesOutDoc"`
						       Cached           int    `json:"cached"`
						       ChromeUserTiming []struct {
							       Name string `json:"name"`
							       Time int    `json:"time"`
						       } `json:"chromeUserTiming"`
						       ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int `json:"connections"`
						       ConsoleLog                                  []struct {
							       Column             int    `json:"column"`
							       ExecutionContextID int    `json:"executionContextId"`
							       Level              string `json:"level"`
							       Line               int    `json:"line"`

							       Source string `json:"source"`

							       Text      string  `json:"text"`
							       Timestamp float64 `json:"timestamp"`
							       Type      string  `json:"type"`
							       URL       string  `json:"url"`
						       } `json:"consoleLog"`
						       Date                       int     `json:"date"`
						       DocCPUms                   float64 `json:"docCPUms"`
						       DocCPUpct                  int     `json:"docCPUpct"`
						       DocTime                    int     `json:"docTime"`
						       DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
						       DomElements                int     `json:"domElements"`
						       DomTime                    int     `json:"domTime"`
						       Domains                    map[string]struct {
											  Bytes       int `json:"bytes"`
											  Connections int `json:"connections"`
											  Requests    int `json:"requests"`
						       				   } `json:"domains"`
						       EffectiveBps      int     `json:"effectiveBps"`
						       EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
						       FirstPaint        int     `json:"firstPaint"`
						       FixedViewport     int     `json:"fixed_viewport"`
						       FullyLoaded       int     `json:"fullyLoaded"`
						       FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
						       GzipSavings       int     `json:"gzip_savings"`
						       GzipTotal         int     `json:"gzip_total"`
						       ImageSavings      int     `json:"image_savings"`
						       ImageTotal        int     `json:"image_total"`

						       IsResponsive        int    `json:"isResponsive"`
						       LastVisualChange    int    `json:"lastVisualChange"`
						       LoadEventEnd        int    `json:"loadEventEnd"`
						       LoadEventStart      int    `json:"loadEventStart"`
						       LoadTime            int    `json:"loadTime"`
						       MinifySavings       int    `json:"minify_savings"`
						       MinifyTotal         int    `json:"minify_total"`
						       OptimizationChecked int    `json:"optimization_checked"`
						       PageSpeedVersion    string `json:"pageSpeedVersion"`

						       Render   int `json:"render"`

						       RequestsDoc          int    `json:"requestsDoc"`
						       RequestsFull         int    `json:"requestsFull"`
						       Responses200         int    `json:"responses_200"`
						       Responses404         int    `json:"responses_404"`
						       ResponsesOther       int    `json:"responses_other"`
						       Result               int    `json:"result"`
						       Run                  int    `json:"run"`
						       ScoreCache           int    `json:"score_cache"`
						       ScoreCdn             int    `json:"score_cdn"`
						       ScoreCombine         int    `json:"score_combine"`
						       ScoreCompress        int    `json:"score_compress"`
						       ScoreCookies         int    `json:"score_cookies"`
						       ScoreEtags           int    `json:"score_etags"`
						       ScoreGzip            int    `json:"score_gzip"`
						       ScoreKeep_alive      int    `json:"score_keep-alive"`
						       ScoreMinify          int    `json:"score_minify"`
						       ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
						       ServerCount          int    `json:"server_count"`
						       ServerRtt            int    `json:"server_rtt"`
						       Tester               string `json:"tester"`

						       Title       string `json:"title"`
						       TitleTime   int    `json:"titleTime"`
						       VideoFrames []struct {
							       VisuallyComplete int    `json:"VisuallyComplete"`

							       Time             int    `json:"time"`
						       } `json:"videoFrames"`
						       VisualComplete int `json:"visualComplete"`
					       } `json:"firstView"`
			     } `json:"median"`
		     Mobile int    `json:"mobile"`
		     Plr    string `json:"plr"`
		     Runs   struct {
				     One struct {
						 FirstView struct {
								   SpeedIndex  int    `json:"SpeedIndex"`
								   TTFB        int    `json:"TTFB"`
								   URL         string `json:"URL"`
								   AdultSite   int    `json:"adult_site"`
								   Aft         int    `json:"aft"`
								   BasePageCdn string `json:"base_page_cdn"`
								   Breakdown   interface{} `json:"breakdown"`
								   BrowserName      string `json:"browser_name"`
								   BrowserVersion   string `json:"browser_version"`
								   BytesIn          int    `json:"bytesIn"`
								   BytesInDoc       int    `json:"bytesInDoc"`
								   BytesOut         int    `json:"bytesOut"`
								   BytesOutDoc      int    `json:"bytesOutDoc"`
								   Cached           int    `json:"cached"`
								   ChromeUserTiming []struct {
									   Name string `json:"name"`
									   Time int    `json:"time"`
								   } `json:"chromeUserTiming"`
								   ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
								   ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
								   ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
								   ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
								   ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
								   ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
								   ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
								   ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
								   ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
								   ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
								   ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
								   ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
								   ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
								   ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
								   ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
								   Connections                                 int `json:"connections"`
								   ConsoleLog                                  []struct {
									   Column             int    `json:"column"`
									   ExecutionContextID int    `json:"executionContextId"`

									   Line               int    `json:"line"`
									   Source string `json:"source"`
									   Text      string  `json:"text"`
									   Timestamp float64 `json:"timestamp"`

								   } `json:"consoleLog"`
								   Date                       int     `json:"date"`
								   DocCPUms                   float64 `json:"docCPUms"`
								   DocCPUpct                  int     `json:"docCPUpct"`
								   DocTime                    int     `json:"docTime"`
								   DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
								   DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
								   DomElements                int     `json:"domElements"`
								   DomTime                    int     `json:"domTime"`
								   Domains                    map[string]struct {
												   Bytes       int `json:"bytes"`
												   Connections int `json:"connections"`
												   Requests    int `json:"requests"`
								   			       } `json:"domains"`
								   EffectiveBps      int     `json:"effectiveBps"`
								   EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
								   FirstPaint        int     `json:"firstPaint"`
								   FixedViewport     int     `json:"fixed_viewport"`
								   FullyLoaded       int     `json:"fullyLoaded"`
								   FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
								   FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
								   GzipSavings       int     `json:"gzip_savings"`
								   GzipTotal         int     `json:"gzip_total"`
								   ImageSavings      int     `json:"image_savings"`
								   ImageTotal        int     `json:"image_total"`

								   IsResponsive        int    `json:"isResponsive"`
								   LastVisualChange    int    `json:"lastVisualChange"`
								   LoadEventEnd        int    `json:"loadEventEnd"`
								   LoadEventStart      int    `json:"loadEventStart"`
								   LoadTime            int    `json:"loadTime"`
								   MinifySavings       int    `json:"minify_savings"`
								   MinifyTotal         int    `json:"minify_total"`
								   OptimizationChecked int    `json:"optimization_checked"`
								   PageSpeedVersion    string `json:"pageSpeedVersion"`

								   Render   int `json:"render"`
								   Requests []struct {

									   Host string `json:"host"`
									   URL string `json:"url"`
									   ResponseCode string `json:"responseCode"`
									   LoadMs string `json:"load_ms"`
									   TtfbMs interface{} `json:"ttfb_ms"`
									   DnsMs  interface{} `json:"dns_ms"`
									   ConnectMs interface{} `json:"connect_ms"`
									   SslMs  interface{} `json:"ssl_ms"`
									   DownloadMs interface{} `json:"download_ms"`
									   LoadStart string `json:"load_start"`
									   BytesOut string `json:"bytesOut"`
									   BytesIn string `json:"bytesIn"`
									   ObjectSize string `json:"objectSize"`
									   Expires string `json:"expires,omitempty"`
									   CacheControl string `json:"cacheControl,omitempty"`
									   ContentType string `json:"contentType,omitempty"`
									   ContentEncoding string `json:"contentEncoding,omitempty"`
									   Type string `json:"type"`
									   Socket string `json:"socket"`
									   ScoreCache string `json:"score_cache"`
									   ScoreCdn string `json:"score_cdn"`
									   ScoreGzip string `json:"score_gzip"`
									   ScoreCookies string `json:"score_cookies"`
									   ScoreKeepAlive string `json:"score_keep-alive"`
									   ScoreMinify string `json:"score_minify"`
									   ScoreCombine string `json:"score_combine"`
									   ScoreCompress string `json:"score_compress"`
									   ScoreEtags string `json:"score_etags"`
									   IsSecure string `json:"is_secure"`

									   TtfbStart string `json:"ttfb_start"`
									   AllStart string `json:"all_start"`

								   } `json:"requests"`
								   RequestsDoc          int    `json:"requestsDoc"`
								   RequestsFull         int    `json:"requestsFull"`
								   Responses200         int    `json:"responses_200"`
								   Responses404         int    `json:"responses_404"`
								   ResponsesOther       int    `json:"responses_other"`
								   Result               int    `json:"result"`
								   Run                  int    `json:"run"`
								   ScoreCache           int    `json:"score_cache"`
								   ScoreCdn             int    `json:"score_cdn"`
								   ScoreCombine         int    `json:"score_combine"`
								   ScoreCompress        int    `json:"score_compress"`
								   ScoreCookies         int    `json:"score_cookies"`
								   ScoreEtags           int    `json:"score_etags"`
								   ScoreGzip            int    `json:"score_gzip"`
								   ScoreKeep_alive      int    `json:"score_keep-alive"`
								   ScoreMinify          int    `json:"score_minify"`
								   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
								   ServerCount          int    `json:"server_count"`
								   ServerRtt            int    `json:"server_rtt"`
								   Tester               string `json:"tester"`

								   Title       string `json:"title"`
								   TitleTime   int    `json:"titleTime"`

								   VisualComplete int `json:"visualComplete"`
							   } `json:"firstView"`
					 } `json:"1"`
			     } `json:"runs"`

		     SuccessfulFVRuns int    `json:"successfulFVRuns"`
		     URL              string `json:"url"`
	     } `json:"data"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}


//for odd json behavior (empty maps being represented as arrays in domain and breakdown

type Breakdown map[string]struct {
	Bytes    int   `json:"bytes"`
	Color    []int `json:"color"`
	Requests int   `json:"requests"`
}

type JsonResults struct {
	Data struct {
		     Average struct {
				     FirstView struct {
						       SpeedIndex                                  int     `json:"SpeedIndex"`
						       TTFB                                        int     `json:"TTFB"`
						       AdultSite                                   int     `json:"adult_site"`
						       Aft                                         int     `json:"aft"`
						       AvgRun                                      int     `json:"avgRun"`
						       BytesIn                                     int     `json:"bytesIn"`
						       BytesInDoc                                  int     `json:"bytesInDoc"`
						       BytesOut                                    int     `json:"bytesOut"`
						       BytesOutDoc                                 int     `json:"bytesOutDoc"`
						       Cached                                      int     `json:"cached"`
						       ChromeUserTiming_domComplete                int     `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int     `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int     `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int     `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int     `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int     `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int     `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int     `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int     `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int     `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int     `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int     `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int     `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int     `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int     `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int     `json:"connections"`
						       Date                                        int     `json:"date"`
						       DocCPUms                                    int     `json:"docCPUms"`
						       DocCPUpct                                   int     `json:"docCPUpct"`
						       DocTime                                     int     `json:"docTime"`
						       DomContentLoadedEventEnd                    int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart                  int     `json:"domContentLoadedEventStart"`
						       DomElements                                 int     `json:"domElements"`
						       DomTime                                     int     `json:"domTime"`
						       EffectiveBps                                int     `json:"effectiveBps"`
						       EffectiveBpsDoc                             int     `json:"effectiveBpsDoc"`
						       FirstPaint                                  int     `json:"firstPaint"`
						       FixedViewport                               int     `json:"fixed_viewport"`
						       FullyLoaded                                 int     `json:"fullyLoaded"`
						       FullyLoadedCPUms                            float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct                           int     `json:"fullyLoadedCPUpct"`
						       GzipSavings                                 int     `json:"gzip_savings"`
						       GzipTotal                                   int     `json:"gzip_total"`
						       ImageSavings                                int     `json:"image_savings"`
						       ImageTotal                                  int     `json:"image_total"`
						       IsResponsive                                int     `json:"isResponsive"`
						       LastVisualChange                            int     `json:"lastVisualChange"`
						       LoadEventEnd                                int     `json:"loadEventEnd"`
						       LoadEventStart                              int     `json:"loadEventStart"`
						       LoadTime                                    int     `json:"loadTime"`
						       MinifySavings                               int     `json:"minify_savings"`
						       MinifyTotal                                 int     `json:"minify_total"`
						       OptimizationChecked                         int     `json:"optimization_checked"`
						       PageSpeedVersion                            float64 `json:"pageSpeedVersion"`
						       Render                                      int     `json:"render"`
						       Requests                                    int     `json:"requests"`
						       RequestsDoc                                 int     `json:"requestsDoc"`
						       RequestsFull                                int     `json:"requestsFull"`
						       Responses200                                int     `json:"responses_200"`
						       Responses404                                int     `json:"responses_404"`
						       ResponsesOther                              int     `json:"responses_other"`
						       Result                                      int     `json:"result"`
						       Run                                         int     `json:"run"`
						       ScoreCache                                  int     `json:"score_cache"`
						       ScoreCdn                                    int     `json:"score_cdn"`
						       ScoreCombine                                int     `json:"score_combine"`
						       ScoreCompress                               int     `json:"score_compress"`
						       ScoreCookies                                int     `json:"score_cookies"`
						       ScoreEtags                                  int     `json:"score_etags"`
						       ScoreGzip                                   int     `json:"score_gzip"`
						       ScoreKeep_alive                             int     `json:"score_keep-alive"`
						       ScoreMinify                                 int     `json:"score_minify"`
						       ScoreProgressiveJpeg                        int     `json:"score_progressive_jpeg"`
						       ServerCount                                 int     `json:"server_count"`
						       ServerRtt                                   int     `json:"server_rtt"`
						       TitleTime                                   int     `json:"titleTime"`
						       VisualComplete                              int     `json:"visualComplete"`
					       } `json:"firstView"`
				     RepeatView struct {
						       SpeedIndex                                  int     `json:"SpeedIndex"`
						       TTFB                                        int     `json:"TTFB"`
						       AdultSite                                   int     `json:"adult_site"`
						       Aft                                         int     `json:"aft"`
						       AvgRun                                      int     `json:"avgRun"`
						       BytesIn                                     int     `json:"bytesIn"`
						       BytesInDoc                                  int     `json:"bytesInDoc"`
						       BytesOut                                    int     `json:"bytesOut"`
						       BytesOutDoc                                 int     `json:"bytesOutDoc"`
						       Cached                                      int     `json:"cached"`
						       ChromeUserTiming_domComplete                int     `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int     `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int     `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int     `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int     `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int     `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int     `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int     `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int     `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int     `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int     `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int     `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int     `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int     `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int     `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int     `json:"connections"`
						       Date                                        int     `json:"date"`
						       DocCPUms                                    float64 `json:"docCPUms"`
						       DocCPUpct                                   int     `json:"docCPUpct"`
						       DocTime                                     int     `json:"docTime"`
						       DomContentLoadedEventEnd                    int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart                  int     `json:"domContentLoadedEventStart"`
						       DomElements                                 int     `json:"domElements"`
						       DomTime                                     int     `json:"domTime"`
						       EffectiveBps                                int     `json:"effectiveBps"`
						       EffectiveBpsDoc                             int     `json:"effectiveBpsDoc"`
						       FirstPaint                                  int     `json:"firstPaint"`
						       FixedViewport                               int     `json:"fixed_viewport"`
						       FullyLoaded                                 int     `json:"fullyLoaded"`
						       FullyLoadedCPUms                            float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct                           int     `json:"fullyLoadedCPUpct"`
						       GzipSavings                                 int     `json:"gzip_savings"`
						       GzipTotal                                   int     `json:"gzip_total"`
						       ImageSavings                                int     `json:"image_savings"`
						       ImageTotal                                  int     `json:"image_total"`
						       IsResponsive                                int     `json:"isResponsive"`
						       LastVisualChange                            int     `json:"lastVisualChange"`
						       LoadEventEnd                                int     `json:"loadEventEnd"`
						       LoadEventStart                              int     `json:"loadEventStart"`
						       LoadTime                                    int     `json:"loadTime"`
						       MinifySavings                               int     `json:"minify_savings"`
						       MinifyTotal                                 int     `json:"minify_total"`
						       OptimizationChecked                         int     `json:"optimization_checked"`
						       PageSpeedVersion                            float64 `json:"pageSpeedVersion"`
						       Render                                      int     `json:"render"`
						       Requests                                    int     `json:"requests"`
						       RequestsDoc                                 int     `json:"requestsDoc"`
						       RequestsFull                                int     `json:"requestsFull"`
						       Responses200                                int     `json:"responses_200"`
						       Responses404                                int     `json:"responses_404"`
						       ResponsesOther                              int     `json:"responses_other"`
						       Result                                      int     `json:"result"`
						       Run                                         int     `json:"run"`
						       ScoreCache                                  int     `json:"score_cache"`
						       ScoreCdn                                    int     `json:"score_cdn"`
						       ScoreCombine                                int     `json:"score_combine"`
						       ScoreCompress                               int     `json:"score_compress"`
						       ScoreCookies                                int     `json:"score_cookies"`
						       ScoreEtags                                  int     `json:"score_etags"`
						       ScoreGzip                                   int     `json:"score_gzip"`
						       ScoreKeep_alive                             int     `json:"score_keep-alive"`
						       ScoreMinify                                 int     `json:"score_minify"`
						       ScoreProgressiveJpeg                        int     `json:"score_progressive_jpeg"`
						       ServerCount                                 int     `json:"server_count"`
						       ServerRtt                                   int     `json:"server_rtt"`
						       TitleTime                                   int     `json:"titleTime"`
						       VisualComplete                              int     `json:"visualComplete"`
					       } `json:"repeatView"`
			     } `json:"average"`
		     BwDown       int    `json:"bwDown"`
		     BwUp         int    `json:"bwUp"`
		     Completed    int    `json:"completed"`
		     Connectivity string `json:"connectivity"`
		     From         string `json:"from"`
		     Fvonly       bool   `json:"fvonly"`
		     ID           string `json:"id"`
		     Latency      int    `json:"latency"`
		     Location     string `json:"location"`
		     Median       struct {
				     FirstView struct {
						       SpeedIndex  int    `json:"SpeedIndex"`
						       TTFB        int    `json:"TTFB"`
						       URL         string `json:"URL"`
						       AdultSite   int    `json:"adult_site"`
						       Aft         int    `json:"aft"`
						       BasePageCdn string `json:"base_page_cdn"`
						       Breakdown   struct {
									   CSS struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"css"`
									   Flash struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"flash"`
									   Font struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"font"`
									   HTML struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"html"`
									   Image struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"image"`
									   Js struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"js"`
									   Other struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"other"`
								   } `json:"breakdown"`
						       BrowserName      string `json:"browser_name"`
						       BrowserVersion   string `json:"browser_version"`
						       BytesIn          int    `json:"bytesIn"`
						       BytesInDoc       int    `json:"bytesInDoc"`
						       BytesOut         int    `json:"bytesOut"`
						       BytesOutDoc      int    `json:"bytesOutDoc"`
						       Cached           int    `json:"cached"`
						       ChromeUserTiming []struct {
							       Name string `json:"name"`
							       Time int    `json:"time"`
						       } `json:"chromeUserTiming"`
						       ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int `json:"connections"`
						       ConsoleLog                                  []struct {
							       Column             int    `json:"column"`
							       ExecutionContextID int    `json:"executionContextId"`
							       Level              string `json:"level"`
							       Line               int    `json:"line"`
							       Parameters         []struct {
								       Type  string `json:"type"`
								       Value string `json:"value"`
							       } `json:"parameters"`
							       Source string `json:"source"`
							       Stack  struct {
											  CallFrames []struct {
												  ColumnNumber int    `json:"columnNumber"`
												  FunctionName string `json:"functionName"`
												  LineNumber   int    `json:"lineNumber"`
												  ScriptID     string `json:"scriptId"`
												  URL          string `json:"url"`
											  } `json:"callFrames"`
										  } `json:"stack"`
							       Text      string  `json:"text"`
							       Timestamp float64 `json:"timestamp"`
							       Type      string  `json:"type"`
							       URL       string  `json:"url"`
						       } `json:"consoleLog"`
						       Date                       int `json:"date"`
						       DocCPUms                   int `json:"docCPUms"`
						       DocCPUpct                  int `json:"docCPUpct"`
						       DocTime                    int `json:"docTime"`
						       DomContentLoadedEventEnd   int `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart int `json:"domContentLoadedEventStart"`
						       DomElements                int `json:"domElements"`
						       DomTime                    int `json:"domTime"`
						       Domains                    struct {
									   UI_portalserver_int_cisco_com struct {
														 Bytes       int `json:"bytes"`
														 Connections int `json:"connections"`
														 Requests    int `json:"requests"`
													 } `json:"ui-portalserver.int.cisco.com"`
								   } `json:"domains"`
						       EffectiveBps      int     `json:"effectiveBps"`
						       EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
						       FirstPaint        int     `json:"firstPaint"`
						       FixedViewport     int     `json:"fixed_viewport"`
						       FullyLoaded       int     `json:"fullyLoaded"`
						       FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
						       GzipSavings       int     `json:"gzip_savings"`
						       GzipTotal         int     `json:"gzip_total"`
						       ImageSavings      int     `json:"image_savings"`
						       ImageTotal        int     `json:"image_total"`
						       Images            struct {
									   Checklist      string `json:"checklist"`
									   ConnectionView string `json:"connectionView"`
									   ScreenShot     string `json:"screenShot"`
									   Waterfall      string `json:"waterfall"`
								   } `json:"images"`
						       IsResponsive        int    `json:"isResponsive"`
						       LastVisualChange    int    `json:"lastVisualChange"`
						       LoadEventEnd        int    `json:"loadEventEnd"`
						       LoadEventStart      int    `json:"loadEventStart"`
						       LoadTime            int    `json:"loadTime"`
						       MinifySavings       int    `json:"minify_savings"`
						       MinifyTotal         int    `json:"minify_total"`
						       OptimizationChecked int    `json:"optimization_checked"`
						       PageSpeedVersion    string `json:"pageSpeedVersion"`
						       Pages               struct {
									   Breakdown  string `json:"breakdown"`
									   Checklist  string `json:"checklist"`
									   Details    string `json:"details"`
									   Domains    string `json:"domains"`
									   ScreenShot string `json:"screenShot"`
								   } `json:"pages"`
						       RawData struct {
									   Headers      string `json:"headers"`
									   PageData     string `json:"pageData"`
									   RequestsData string `json:"requestsData"`
									   Utilization  string `json:"utilization"`
								   } `json:"rawData"`
						       Render   int `json:"render"`
						       Requests []struct {
							       AllEnd        int    `json:"all_end"`
							       AllMs         int    `json:"all_ms"`
							       AllStart      string `json:"all_start"`
							       BytesIn       string `json:"bytesIn"`
							       BytesOut      string `json:"bytesOut"`
							       CacheControl  string `json:"cacheControl"`
							       CacheTime     string `json:"cache_time"`
							       ClientPort    string `json:"client_port"`
							       ConnectEnd    string `json:"connect_end"`
							       ConnectMs     string `json:"connect_ms"`
							       ConnectStart  string `json:"connect_start"`
							       ContentType   string `json:"contentType"`
							       DNSEnd        string `json:"dns_end"`
							       DNSMs         string `json:"dns_ms"`
							       DNSStart      string `json:"dns_start"`
							       DownloadEnd   int    `json:"download_end"`
							       DownloadMs    int    `json:"download_ms"`
							       DownloadStart int    `json:"download_start"`
							       FullURL       string `json:"full_url"`
							       GzipSave      string `json:"gzip_save"`
							       GzipTotal     string `json:"gzip_total"`
							       Headers       struct {
										     Request  []string `json:"request"`
										     Response []string `json:"response"`
									     } `json:"headers"`
							       Host                 string `json:"host"`
							       ImageSave            string `json:"image_save"`
							       ImageTotal           string `json:"image_total"`
							       Index                int    `json:"index"`
							       IsSecure             string `json:"is_secure"`
							       JpegScanCount        string `json:"jpeg_scan_count"`
							       LoadEnd              int    `json:"load_end"`
							       LoadMs               string `json:"load_ms"`
							       LoadStart            string `json:"load_start"`
							       Method               string `json:"method"`
							       MinifySave           string `json:"minify_save"`
							       MinifyTotal          string `json:"minify_total"`
							       Number               int    `json:"number"`
							       ObjectSize           string `json:"objectSize"`
							       Priority             string `json:"priority"`
							       ResponseCode         string `json:"responseCode"`
							       ScoreCache           string `json:"score_cache"`
							       ScoreCdn             string `json:"score_cdn"`
							       ScoreCombine         string `json:"score_combine"`
							       ScoreCompress        string `json:"score_compress"`
							       ScoreCookies         string `json:"score_cookies"`
							       ScoreEtags           string `json:"score_etags"`
							       ScoreGzip            string `json:"score_gzip"`
							       ScoreKeep_alive      string `json:"score_keep-alive"`
							       ScoreMinify          string `json:"score_minify"`
							       ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
							       ServerCount          string `json:"server_count"`
							       Socket               string `json:"socket"`
							       SslEnd               string `json:"ssl_end"`
							       SslMs                int    `json:"ssl_ms"`
							       SslStart             string `json:"ssl_start"`
							       TtfbEnd              int    `json:"ttfb_end"`
							       TtfbMs               string `json:"ttfb_ms"`
							       TtfbStart            string `json:"ttfb_start"`
							       Type                 string `json:"type"`
							       URL                  string `json:"url"`
						       } `json:"requests"`
						       RequestsDoc          int    `json:"requestsDoc"`
						       RequestsFull         int    `json:"requestsFull"`
						       Responses200         int    `json:"responses_200"`
						       Responses404         int    `json:"responses_404"`
						       ResponsesOther       int    `json:"responses_other"`
						       Result               int    `json:"result"`
						       Run                  int    `json:"run"`
						       ScoreCache           int    `json:"score_cache"`
						       ScoreCdn             int    `json:"score_cdn"`
						       ScoreCombine         int    `json:"score_combine"`
						       ScoreCompress        int    `json:"score_compress"`
						       ScoreCookies         int    `json:"score_cookies"`
						       ScoreEtags           int    `json:"score_etags"`
						       ScoreGzip            int    `json:"score_gzip"`
						       ScoreKeep_alive      int    `json:"score_keep-alive"`
						       ScoreMinify          int    `json:"score_minify"`
						       ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
						       ServerCount          int    `json:"server_count"`
						       ServerRtt            int    `json:"server_rtt"`
						       Tester               string `json:"tester"`
						       Thumbnails           struct {
									   Checklist  string `json:"checklist"`
									   ScreenShot string `json:"screenShot"`
									   Waterfall  string `json:"waterfall"`
								   } `json:"thumbnails"`
						       Title       string `json:"title"`
						       TitleTime   int    `json:"titleTime"`
						       VideoFrames []struct {
							       VisuallyComplete int    `json:"VisuallyComplete"`
							       Image            string `json:"image"`
							       Time             int    `json:"time"`
						       } `json:"videoFrames"`
						       VisualComplete int `json:"visualComplete"`
					       } `json:"firstView"`
				     RepeatView struct {
						       SpeedIndex  int    `json:"SpeedIndex"`
						       TTFB        int    `json:"TTFB"`
						       URL         string `json:"URL"`
						       AdultSite   int    `json:"adult_site"`
						       Aft         int    `json:"aft"`
						       BasePageCdn string `json:"base_page_cdn"`
						       Breakdown   struct {
									   CSS struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"css"`
									   Flash struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"flash"`
									   Font struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"font"`
									   HTML struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"html"`
									   Image struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"image"`
									   Js struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"js"`
									   Other struct {
										       Bytes    int   `json:"bytes"`
										       Color    []int `json:"color"`
										       Requests int   `json:"requests"`
									       } `json:"other"`
								   } `json:"breakdown"`
						       BrowserName      string `json:"browser_name"`
						       BrowserVersion   string `json:"browser_version"`
						       BytesIn          int    `json:"bytesIn"`
						       BytesInDoc       int    `json:"bytesInDoc"`
						       BytesOut         int    `json:"bytesOut"`
						       BytesOutDoc      int    `json:"bytesOutDoc"`
						       Cached           int    `json:"cached"`
						       ChromeUserTiming []struct {
							       Name string `json:"name"`
							       Time int    `json:"time"`
						       } `json:"chromeUserTiming"`
						       ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int `json:"connections"`
						       ConsoleLog                                  []struct {
							       Column             int    `json:"column"`
							       ExecutionContextID int    `json:"executionContextId"`
							       Level              string `json:"level"`
							       Line               int    `json:"line"`
							       Parameters         []struct {
								       Type  string `json:"type"`
								       Value string `json:"value"`
							       } `json:"parameters"`
							       Source string `json:"source"`
							       Stack  struct {
											  CallFrames []struct {
												  ColumnNumber int    `json:"columnNumber"`
												  FunctionName string `json:"functionName"`
												  LineNumber   int    `json:"lineNumber"`
												  ScriptID     string `json:"scriptId"`
												  URL          string `json:"url"`
											  } `json:"callFrames"`
										  } `json:"stack"`
							       Text      string  `json:"text"`
							       Timestamp float64 `json:"timestamp"`
							       Type      string  `json:"type"`
							       URL       string  `json:"url"`
						       } `json:"consoleLog"`
						       Date                       int     `json:"date"`
						       DocCPUms                   float64 `json:"docCPUms"`
						       DocCPUpct                  int     `json:"docCPUpct"`
						       DocTime                    int     `json:"docTime"`
						       DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
						       DomElements                int     `json:"domElements"`
						       DomTime                    int     `json:"domTime"`
						       Domains                    struct {
									   UI_portalserver_int_cisco_com struct {
														 Bytes       int `json:"bytes"`
														 Connections int `json:"connections"`
														 Requests    int `json:"requests"`
													 } `json:"ui-portalserver.int.cisco.com"`
								   } `json:"domains"`
						       EffectiveBps      int     `json:"effectiveBps"`
						       EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
						       FirstPaint        int     `json:"firstPaint"`
						       FixedViewport     int     `json:"fixed_viewport"`
						       FullyLoaded       int     `json:"fullyLoaded"`
						       FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
						       GzipSavings       int     `json:"gzip_savings"`
						       GzipTotal         int     `json:"gzip_total"`
						       ImageSavings      int     `json:"image_savings"`
						       ImageTotal        int     `json:"image_total"`
						       Images            struct {
									   Checklist      string `json:"checklist"`
									   ConnectionView string `json:"connectionView"`
									   ScreenShot     string `json:"screenShot"`
									   Waterfall      string `json:"waterfall"`
								   } `json:"images"`
						       IsResponsive        int    `json:"isResponsive"`
						       LastVisualChange    int    `json:"lastVisualChange"`
						       LoadEventEnd        int    `json:"loadEventEnd"`
						       LoadEventStart      int    `json:"loadEventStart"`
						       LoadTime            int    `json:"loadTime"`
						       MinifySavings       int    `json:"minify_savings"`
						       MinifyTotal         int    `json:"minify_total"`
						       OptimizationChecked int    `json:"optimization_checked"`
						       PageSpeedVersion    string `json:"pageSpeedVersion"`
						       Pages               struct {
									   Breakdown  string `json:"breakdown"`
									   Checklist  string `json:"checklist"`
									   Details    string `json:"details"`
									   Domains    string `json:"domains"`
									   ScreenShot string `json:"screenShot"`
								   } `json:"pages"`
						       RawData struct {
									   Headers      string `json:"headers"`
									   PageData     string `json:"pageData"`
									   RequestsData string `json:"requestsData"`
									   Utilization  string `json:"utilization"`
								   } `json:"rawData"`
						       Render   int `json:"render"`
						       Requests []struct {
							       AllEnd        int    `json:"all_end"`
							       AllMs         int    `json:"all_ms"`
							       AllStart      string `json:"all_start"`
							       BytesIn       string `json:"bytesIn"`
							       BytesOut      string `json:"bytesOut"`
							       CacheControl  string `json:"cacheControl"`
							       CacheTime     string `json:"cache_time"`
							       ClientPort    string `json:"client_port"`
							       ConnectEnd    string `json:"connect_end"`
							       ConnectMs     string `json:"connect_ms"`
							       ConnectStart  string `json:"connect_start"`
							       ContentType   string `json:"contentType"`
							       DNSEnd        string `json:"dns_end"`
							       DNSMs         string `json:"dns_ms"`
							       DNSStart      string `json:"dns_start"`
							       DownloadEnd   int    `json:"download_end"`
							       DownloadMs    int    `json:"download_ms"`
							       DownloadStart int    `json:"download_start"`
							       FullURL       string `json:"full_url"`
							       GzipSave      string `json:"gzip_save"`
							       GzipTotal     string `json:"gzip_total"`
							       Headers       struct {
										     Request  []string `json:"request"`
										     Response []string `json:"response"`
									     } `json:"headers"`
							       Host                 string `json:"host"`
							       ImageSave            string `json:"image_save"`
							       ImageTotal           string `json:"image_total"`
							       Index                int    `json:"index"`
							       IsSecure             string `json:"is_secure"`
							       JpegScanCount        string `json:"jpeg_scan_count"`
							       LoadEnd              int    `json:"load_end"`
							       LoadMs               string `json:"load_ms"`
							       LoadStart            string `json:"load_start"`
							       Method               string `json:"method"`
							       MinifySave           string `json:"minify_save"`
							       MinifyTotal          string `json:"minify_total"`
							       Number               int    `json:"number"`
							       ObjectSize           string `json:"objectSize"`
							       Priority             string `json:"priority"`
							       ResponseCode         string `json:"responseCode"`
							       ScoreCache           string `json:"score_cache"`
							       ScoreCdn             string `json:"score_cdn"`
							       ScoreCombine         string `json:"score_combine"`
							       ScoreCompress        string `json:"score_compress"`
							       ScoreCookies         string `json:"score_cookies"`
							       ScoreEtags           string `json:"score_etags"`
							       ScoreGzip            string `json:"score_gzip"`
							       ScoreKeep_alive      string `json:"score_keep-alive"`
							       ScoreMinify          string `json:"score_minify"`
							       ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
							       ServerCount          string `json:"server_count"`
							       Socket               string `json:"socket"`
							       SslEnd               string `json:"ssl_end"`
							       SslMs                int    `json:"ssl_ms"`
							       SslStart             string `json:"ssl_start"`
							       TtfbEnd              int    `json:"ttfb_end"`
							       TtfbMs               string `json:"ttfb_ms"`
							       TtfbStart            string `json:"ttfb_start"`
							       Type                 string `json:"type"`
							       URL                  string `json:"url"`
						       } `json:"requests"`
						       RequestsDoc          int    `json:"requestsDoc"`
						       RequestsFull         int    `json:"requestsFull"`
						       Responses200         int    `json:"responses_200"`
						       Responses404         int    `json:"responses_404"`
						       ResponsesOther       int    `json:"responses_other"`
						       Result               int    `json:"result"`
						       Run                  int    `json:"run"`
						       ScoreCache           int    `json:"score_cache"`
						       ScoreCdn             int    `json:"score_cdn"`
						       ScoreCombine         int    `json:"score_combine"`
						       ScoreCompress        int    `json:"score_compress"`
						       ScoreCookies         int    `json:"score_cookies"`
						       ScoreEtags           int    `json:"score_etags"`
						       ScoreGzip            int    `json:"score_gzip"`
						       ScoreKeep_alive      int    `json:"score_keep-alive"`
						       ScoreMinify          int    `json:"score_minify"`
						       ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
						       ServerCount          int    `json:"server_count"`
						       ServerRtt            int    `json:"server_rtt"`
						       Tester               string `json:"tester"`
						       Thumbnails           struct {
									   Checklist  string `json:"checklist"`
									   ScreenShot string `json:"screenShot"`
									   Waterfall  string `json:"waterfall"`
								   } `json:"thumbnails"`
						       Title       string `json:"title"`
						       TitleTime   int    `json:"titleTime"`
						       VideoFrames []struct {
							       VisuallyComplete int    `json:"VisuallyComplete"`
							       Image            string `json:"image"`
							       Time             int    `json:"time"`
						       } `json:"videoFrames"`
						       VisualComplete int `json:"visualComplete"`
					       } `json:"repeatView"`
			     } `json:"median"`
		     Mobile int    `json:"mobile"`
		     Plr    string `json:"plr"`
		     Runs   struct {
				     One struct {
						 FirstView struct {
								   SpeedIndex  int    `json:"SpeedIndex"`
								   TTFB        int    `json:"TTFB"`
								   URL         string `json:"URL"`
								   AdultSite   int    `json:"adult_site"`
								   Aft         int    `json:"aft"`
								   BasePageCdn string `json:"base_page_cdn"`
								   Breakdown   struct {
										       CSS struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"css"`
										       Flash struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"flash"`
										       Font struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"font"`
										       HTML struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"html"`
										       Image struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"image"`
										       Js struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"js"`
										       Other struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"other"`
									       } `json:"breakdown"`
								   BrowserName      string `json:"browser_name"`
								   BrowserVersion   string `json:"browser_version"`
								   BytesIn          int    `json:"bytesIn"`
								   BytesInDoc       int    `json:"bytesInDoc"`
								   BytesOut         int    `json:"bytesOut"`
								   BytesOutDoc      int    `json:"bytesOutDoc"`
								   Cached           int    `json:"cached"`
								   ChromeUserTiming []struct {
									   Name string `json:"name"`
									   Time int    `json:"time"`
								   } `json:"chromeUserTiming"`
								   ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
								   ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
								   ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
								   ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
								   ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
								   ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
								   ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
								   ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
								   ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
								   ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
								   ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
								   ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
								   ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
								   ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
								   ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
								   Connections                                 int `json:"connections"`
								   ConsoleLog                                  []struct {
									   Column             int    `json:"column"`
									   ExecutionContextID int    `json:"executionContextId"`
									   Level              string `json:"level"`
									   Line               int    `json:"line"`
									   Parameters         []struct {
										   Type  string `json:"type"`
										   Value string `json:"value"`
									   } `json:"parameters"`
									   Source string `json:"source"`
									   Stack  struct {
												      CallFrames []struct {
													      ColumnNumber int    `json:"columnNumber"`
													      FunctionName string `json:"functionName"`
													      LineNumber   int    `json:"lineNumber"`
													      ScriptID     string `json:"scriptId"`
													      URL          string `json:"url"`
												      } `json:"callFrames"`
											      } `json:"stack"`
									   Text      string  `json:"text"`
									   Timestamp float64 `json:"timestamp"`
									   Type      string  `json:"type"`
									   URL       string  `json:"url"`
								   } `json:"consoleLog"`
								   Date                       int `json:"date"`
								   DocCPUms                   int `json:"docCPUms"`
								   DocCPUpct                  int `json:"docCPUpct"`
								   DocTime                    int `json:"docTime"`
								   DomContentLoadedEventEnd   int `json:"domContentLoadedEventEnd"`
								   DomContentLoadedEventStart int `json:"domContentLoadedEventStart"`
								   DomElements                int `json:"domElements"`
								   DomTime                    int `json:"domTime"`
								   Domains                    struct {
										       UI_portalserver_int_cisco_com struct {
															     Bytes       int `json:"bytes"`
															     Connections int `json:"connections"`
															     Requests    int `json:"requests"`
														     } `json:"ui-portalserver.int.cisco.com"`
									       } `json:"domains"`
								   EffectiveBps      int     `json:"effectiveBps"`
								   EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
								   FirstPaint        int     `json:"firstPaint"`
								   FixedViewport     int     `json:"fixed_viewport"`
								   FullyLoaded       int     `json:"fullyLoaded"`
								   FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
								   FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
								   GzipSavings       int     `json:"gzip_savings"`
								   GzipTotal         int     `json:"gzip_total"`
								   ImageSavings      int     `json:"image_savings"`
								   ImageTotal        int     `json:"image_total"`
								   Images            struct {
										       Checklist      string `json:"checklist"`
										       ConnectionView string `json:"connectionView"`
										       ScreenShot     string `json:"screenShot"`
										       Waterfall      string `json:"waterfall"`
									       } `json:"images"`
								   IsResponsive        int    `json:"isResponsive"`
								   LastVisualChange    int    `json:"lastVisualChange"`
								   LoadEventEnd        int    `json:"loadEventEnd"`
								   LoadEventStart      int    `json:"loadEventStart"`
								   LoadTime            int    `json:"loadTime"`
								   MinifySavings       int    `json:"minify_savings"`
								   MinifyTotal         int    `json:"minify_total"`
								   OptimizationChecked int    `json:"optimization_checked"`
								   PageSpeedVersion    string `json:"pageSpeedVersion"`
								   Pages               struct {
										       Breakdown  string `json:"breakdown"`
										       Checklist  string `json:"checklist"`
										       Details    string `json:"details"`
										       Domains    string `json:"domains"`
										       ScreenShot string `json:"screenShot"`
									       } `json:"pages"`
								   RawData struct {
										       Headers      string `json:"headers"`
										       PageData     string `json:"pageData"`
										       RequestsData string `json:"requestsData"`
										       Utilization  string `json:"utilization"`
									       } `json:"rawData"`
								   Render   int `json:"render"`
								   Requests []struct {
									   AllEnd        int    `json:"all_end"`
									   AllMs         int    `json:"all_ms"`
									   AllStart      string `json:"all_start"`
									   BytesIn       string `json:"bytesIn"`
									   BytesOut      string `json:"bytesOut"`
									   CacheControl  string `json:"cacheControl"`
									   CacheTime     string `json:"cache_time"`
									   ClientPort    string `json:"client_port"`
									   ConnectEnd    string `json:"connect_end"`
									   ConnectMs     string `json:"connect_ms"`
									   ConnectStart  string `json:"connect_start"`
									   ContentType   string `json:"contentType"`
									   DNSEnd        string `json:"dns_end"`
									   DNSMs         string `json:"dns_ms"`
									   DNSStart      string `json:"dns_start"`
									   DownloadEnd   int    `json:"download_end"`
									   DownloadMs    int    `json:"download_ms"`
									   DownloadStart int    `json:"download_start"`
									   FullURL       string `json:"full_url"`
									   GzipSave      string `json:"gzip_save"`
									   GzipTotal     string `json:"gzip_total"`
									   Headers       struct {
												 Request  []string `json:"request"`
												 Response []string `json:"response"`
											 } `json:"headers"`
									   Host                 string `json:"host"`
									   ImageSave            string `json:"image_save"`
									   ImageTotal           string `json:"image_total"`
									   Index                int    `json:"index"`
									   IsSecure             string `json:"is_secure"`
									   JpegScanCount        string `json:"jpeg_scan_count"`
									   LoadEnd              int    `json:"load_end"`
									   LoadMs               string `json:"load_ms"`
									   LoadStart            string `json:"load_start"`
									   Method               string `json:"method"`
									   MinifySave           string `json:"minify_save"`
									   MinifyTotal          string `json:"minify_total"`
									   Number               int    `json:"number"`
									   ObjectSize           string `json:"objectSize"`
									   Priority             string `json:"priority"`
									   ResponseCode         string `json:"responseCode"`
									   ScoreCache           string `json:"score_cache"`
									   ScoreCdn             string `json:"score_cdn"`
									   ScoreCombine         string `json:"score_combine"`
									   ScoreCompress        string `json:"score_compress"`
									   ScoreCookies         string `json:"score_cookies"`
									   ScoreEtags           string `json:"score_etags"`
									   ScoreGzip            string `json:"score_gzip"`
									   ScoreKeep_alive      string `json:"score_keep-alive"`
									   ScoreMinify          string `json:"score_minify"`
									   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
									   ServerCount          string `json:"server_count"`
									   Socket               string `json:"socket"`
									   SslEnd               string `json:"ssl_end"`
									   SslMs                int    `json:"ssl_ms"`
									   SslStart             string `json:"ssl_start"`
									   TtfbEnd              int    `json:"ttfb_end"`
									   TtfbMs               string `json:"ttfb_ms"`
									   TtfbStart            string `json:"ttfb_start"`
									   Type                 string `json:"type"`
									   URL                  string `json:"url"`
								   } `json:"requests"`
								   RequestsDoc          int    `json:"requestsDoc"`
								   RequestsFull         int    `json:"requestsFull"`
								   Responses200         int    `json:"responses_200"`
								   Responses404         int    `json:"responses_404"`
								   ResponsesOther       int    `json:"responses_other"`
								   Result               int    `json:"result"`
								   Run                  int    `json:"run"`
								   ScoreCache           int    `json:"score_cache"`
								   ScoreCdn             int    `json:"score_cdn"`
								   ScoreCombine         int    `json:"score_combine"`
								   ScoreCompress        int    `json:"score_compress"`
								   ScoreCookies         int    `json:"score_cookies"`
								   ScoreEtags           int    `json:"score_etags"`
								   ScoreGzip            int    `json:"score_gzip"`
								   ScoreKeep_alive      int    `json:"score_keep-alive"`
								   ScoreMinify          int    `json:"score_minify"`
								   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
								   ServerCount          int    `json:"server_count"`
								   ServerRtt            int    `json:"server_rtt"`
								   Tester               string `json:"tester"`
								   Thumbnails           struct {
										       Checklist  string `json:"checklist"`
										       ScreenShot string `json:"screenShot"`
										       Waterfall  string `json:"waterfall"`
									       } `json:"thumbnails"`
								   Title       string `json:"title"`
								   TitleTime   int    `json:"titleTime"`
								   VideoFrames []struct {
									   VisuallyComplete int    `json:"VisuallyComplete"`
									   Image            string `json:"image"`
									   Time             int    `json:"time"`
								   } `json:"videoFrames"`
								   VisualComplete int `json:"visualComplete"`
							   } `json:"firstView"`
						 RepeatView struct {
								   SpeedIndex  int    `json:"SpeedIndex"`
								   TTFB        int    `json:"TTFB"`
								   URL         string `json:"URL"`
								   AdultSite   int    `json:"adult_site"`
								   Aft         int    `json:"aft"`
								   BasePageCdn string `json:"base_page_cdn"`
								   Breakdown   struct {
										       CSS struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"css"`
										       Flash struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"flash"`
										       Font struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"font"`
										       HTML struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"html"`
										       Image struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"image"`
										       Js struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"js"`
										       Other struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"other"`
									       } `json:"breakdown"`
								   BrowserName      string `json:"browser_name"`
								   BrowserVersion   string `json:"browser_version"`
								   BytesIn          int    `json:"bytesIn"`
								   BytesInDoc       int    `json:"bytesInDoc"`
								   BytesOut         int    `json:"bytesOut"`
								   BytesOutDoc      int    `json:"bytesOutDoc"`
								   Cached           int    `json:"cached"`
								   ChromeUserTiming []struct {
									   Name string `json:"name"`
									   Time int    `json:"time"`
								   } `json:"chromeUserTiming"`
								   ChromeUserTiming_domComplete                int `json:"chromeUserTiming.domComplete"`
								   ChromeUserTiming_domContentLoadedEventEnd   int `json:"chromeUserTiming.domContentLoadedEventEnd"`
								   ChromeUserTiming_domContentLoadedEventStart int `json:"chromeUserTiming.domContentLoadedEventStart"`
								   ChromeUserTiming_domInteractive             int `json:"chromeUserTiming.domInteractive"`
								   ChromeUserTiming_domLoading                 int `json:"chromeUserTiming.domLoading"`
								   ChromeUserTiming_fetchStart                 int `json:"chromeUserTiming.fetchStart"`
								   ChromeUserTiming_firstContentfulPaint       int `json:"chromeUserTiming.firstContentfulPaint"`
								   ChromeUserTiming_firstLayout                int `json:"chromeUserTiming.firstLayout"`
								   ChromeUserTiming_firstPaint                 int `json:"chromeUserTiming.firstPaint"`
								   ChromeUserTiming_firstTextPaint             int `json:"chromeUserTiming.firstTextPaint"`
								   ChromeUserTiming_loadEventEnd               int `json:"chromeUserTiming.loadEventEnd"`
								   ChromeUserTiming_loadEventStart             int `json:"chromeUserTiming.loadEventStart"`
								   ChromeUserTiming_responseEnd                int `json:"chromeUserTiming.responseEnd"`
								   ChromeUserTiming_unloadEventEnd             int `json:"chromeUserTiming.unloadEventEnd"`
								   ChromeUserTiming_unloadEventStart           int `json:"chromeUserTiming.unloadEventStart"`
								   Connections                                 int `json:"connections"`
								   ConsoleLog                                  []struct {
									   Column             int    `json:"column"`
									   ExecutionContextID int    `json:"executionContextId"`
									   Level              string `json:"level"`
									   Line               int    `json:"line"`
									   Parameters         []struct {
										   Type  string `json:"type"`
										   Value string `json:"value"`
									   } `json:"parameters"`
									   Source string `json:"source"`
									   Stack  struct {
												      CallFrames []struct {
													      ColumnNumber int    `json:"columnNumber"`
													      FunctionName string `json:"functionName"`
													      LineNumber   int    `json:"lineNumber"`
													      ScriptID     string `json:"scriptId"`
													      URL          string `json:"url"`
												      } `json:"callFrames"`
											      } `json:"stack"`
									   Text      string  `json:"text"`
									   Timestamp float64 `json:"timestamp"`
									   Type      string  `json:"type"`
									   URL       string  `json:"url"`
								   } `json:"consoleLog"`
								   Date                       int     `json:"date"`
								   DocCPUms                   float64 `json:"docCPUms"`
								   DocCPUpct                  int     `json:"docCPUpct"`
								   DocTime                    int     `json:"docTime"`
								   DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
								   DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
								   DomElements                int     `json:"domElements"`
								   DomTime                    int     `json:"domTime"`
								   Domains                    struct {
										       UI_portalserver_int_cisco_com struct {
															     Bytes       int `json:"bytes"`
															     Connections int `json:"connections"`
															     Requests    int `json:"requests"`
														     } `json:"ui-portalserver.int.cisco.com"`
									       } `json:"domains"`
								   EffectiveBps      int     `json:"effectiveBps"`
								   EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
								   FirstPaint        int     `json:"firstPaint"`
								   FixedViewport     int     `json:"fixed_viewport"`
								   FullyLoaded       int     `json:"fullyLoaded"`
								   FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
								   FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
								   GzipSavings       int     `json:"gzip_savings"`
								   GzipTotal         int     `json:"gzip_total"`
								   ImageSavings      int     `json:"image_savings"`
								   ImageTotal        int     `json:"image_total"`
								   Images            struct {
										       Checklist      string `json:"checklist"`
										       ConnectionView string `json:"connectionView"`
										       ScreenShot     string `json:"screenShot"`
										       Waterfall      string `json:"waterfall"`
									       } `json:"images"`
								   IsResponsive        int    `json:"isResponsive"`
								   LastVisualChange    int    `json:"lastVisualChange"`
								   LoadEventEnd        int    `json:"loadEventEnd"`
								   LoadEventStart      int    `json:"loadEventStart"`
								   LoadTime            int    `json:"loadTime"`
								   MinifySavings       int    `json:"minify_savings"`
								   MinifyTotal         int    `json:"minify_total"`
								   OptimizationChecked int    `json:"optimization_checked"`
								   PageSpeedVersion    string `json:"pageSpeedVersion"`
								   Pages               struct {
										       Breakdown  string `json:"breakdown"`
										       Checklist  string `json:"checklist"`
										       Details    string `json:"details"`
										       Domains    string `json:"domains"`
										       ScreenShot string `json:"screenShot"`
									       } `json:"pages"`
								   RawData struct {
										       Headers      string `json:"headers"`
										       PageData     string `json:"pageData"`
										       RequestsData string `json:"requestsData"`
										       Utilization  string `json:"utilization"`
									       } `json:"rawData"`
								   Render   int `json:"render"`
								   Requests []struct {
									   AllEnd        int    `json:"all_end"`
									   AllMs         int    `json:"all_ms"`
									   AllStart      string `json:"all_start"`
									   BytesIn       string `json:"bytesIn"`
									   BytesOut      string `json:"bytesOut"`
									   CacheControl  string `json:"cacheControl"`
									   CacheTime     string `json:"cache_time"`
									   ClientPort    string `json:"client_port"`
									   ConnectEnd    string `json:"connect_end"`
									   ConnectMs     string `json:"connect_ms"`
									   ConnectStart  string `json:"connect_start"`
									   ContentType   string `json:"contentType"`
									   DNSEnd        string `json:"dns_end"`
									   DNSMs         string `json:"dns_ms"`
									   DNSStart      string `json:"dns_start"`
									   DownloadEnd   int    `json:"download_end"`
									   DownloadMs    int    `json:"download_ms"`
									   DownloadStart int    `json:"download_start"`
									   FullURL       string `json:"full_url"`
									   GzipSave      string `json:"gzip_save"`
									   GzipTotal     string `json:"gzip_total"`
									   Headers       struct {
												 Request  []string `json:"request"`
												 Response []string `json:"response"`
											 } `json:"headers"`
									   Host                 string `json:"host"`
									   ImageSave            string `json:"image_save"`
									   ImageTotal           string `json:"image_total"`
									   Index                int    `json:"index"`
									   IsSecure             string `json:"is_secure"`
									   JpegScanCount        string `json:"jpeg_scan_count"`
									   LoadEnd              int    `json:"load_end"`
									   LoadMs               string `json:"load_ms"`
									   LoadStart            string `json:"load_start"`
									   Method               string `json:"method"`
									   MinifySave           string `json:"minify_save"`
									   MinifyTotal          string `json:"minify_total"`
									   Number               int    `json:"number"`
									   ObjectSize           string `json:"objectSize"`
									   Priority             string `json:"priority"`
									   ResponseCode         string `json:"responseCode"`
									   ScoreCache           string `json:"score_cache"`
									   ScoreCdn             string `json:"score_cdn"`
									   ScoreCombine         string `json:"score_combine"`
									   ScoreCompress        string `json:"score_compress"`
									   ScoreCookies         string `json:"score_cookies"`
									   ScoreEtags           string `json:"score_etags"`
									   ScoreGzip            string `json:"score_gzip"`
									   ScoreKeep_alive      string `json:"score_keep-alive"`
									   ScoreMinify          string `json:"score_minify"`
									   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
									   ServerCount          string `json:"server_count"`
									   Socket               string `json:"socket"`
									   SslEnd               string `json:"ssl_end"`
									   SslMs                int    `json:"ssl_ms"`
									   SslStart             string `json:"ssl_start"`
									   TtfbEnd              int    `json:"ttfb_end"`
									   TtfbMs               string `json:"ttfb_ms"`
									   TtfbStart            string `json:"ttfb_start"`
									   Type                 string `json:"type"`
									   URL                  string `json:"url"`
								   } `json:"requests"`
								   RequestsDoc          int    `json:"requestsDoc"`
								   RequestsFull         int    `json:"requestsFull"`
								   Responses200         int    `json:"responses_200"`
								   Responses404         int    `json:"responses_404"`
								   ResponsesOther       int    `json:"responses_other"`
								   Result               int    `json:"result"`
								   Run                  int    `json:"run"`
								   ScoreCache           int    `json:"score_cache"`
								   ScoreCdn             int    `json:"score_cdn"`
								   ScoreCombine         int    `json:"score_combine"`
								   ScoreCompress        int    `json:"score_compress"`
								   ScoreCookies         int    `json:"score_cookies"`
								   ScoreEtags           int    `json:"score_etags"`
								   ScoreGzip            int    `json:"score_gzip"`
								   ScoreKeep_alive      int    `json:"score_keep-alive"`
								   ScoreMinify          int    `json:"score_minify"`
								   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
								   ServerCount          int    `json:"server_count"`
								   ServerRtt            int    `json:"server_rtt"`
								   Tester               string `json:"tester"`
								   Thumbnails           struct {
										       Checklist  string `json:"checklist"`
										       ScreenShot string `json:"screenShot"`
										       Waterfall  string `json:"waterfall"`
									       } `json:"thumbnails"`
								   Title       string `json:"title"`
								   TitleTime   int    `json:"titleTime"`
								   VideoFrames []struct {
									   VisuallyComplete int    `json:"VisuallyComplete"`
									   Image            string `json:"image"`
									   Time             int    `json:"time"`
								   } `json:"videoFrames"`
								   VisualComplete int `json:"visualComplete"`
							   } `json:"repeatView"`
					 } `json:"1"`
			     } `json:"runs"`
		     StandardDeviation struct {
				     FirstView struct {
						       SpeedIndex                                  int         `json:"SpeedIndex"`
						       TTFB                                        int         `json:"TTFB"`
						       AdultSite                                   int         `json:"adult_site"`
						       Aft                                         int         `json:"aft"`
						       AvgRun                                      interface{} `json:"avgRun"`
						       BytesIn                                     int         `json:"bytesIn"`
						       BytesInDoc                                  int         `json:"bytesInDoc"`
						       BytesOut                                    int         `json:"bytesOut"`
						       BytesOutDoc                                 int         `json:"bytesOutDoc"`
						       Cached                                      int         `json:"cached"`
						       ChromeUserTiming_domComplete                int         `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int         `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int         `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int         `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int         `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int         `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int         `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int         `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int         `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int         `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int         `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int         `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int         `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int         `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int         `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int         `json:"connections"`
						       Date                                        int         `json:"date"`
						       DocCPUms                                    int         `json:"docCPUms"`
						       DocCPUpct                                   int         `json:"docCPUpct"`
						       DocTime                                     int         `json:"docTime"`
						       DomContentLoadedEventEnd                    int         `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart                  int         `json:"domContentLoadedEventStart"`
						       DomElements                                 int         `json:"domElements"`
						       DomTime                                     int         `json:"domTime"`
						       EffectiveBps                                int         `json:"effectiveBps"`
						       EffectiveBpsDoc                             int         `json:"effectiveBpsDoc"`
						       FirstPaint                                  int         `json:"firstPaint"`
						       FixedViewport                               int         `json:"fixed_viewport"`
						       FullyLoaded                                 int         `json:"fullyLoaded"`
						       FullyLoadedCPUms                            int         `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct                           int         `json:"fullyLoadedCPUpct"`
						       GzipSavings                                 int         `json:"gzip_savings"`
						       GzipTotal                                   int         `json:"gzip_total"`
						       ImageSavings                                int         `json:"image_savings"`
						       ImageTotal                                  int         `json:"image_total"`
						       IsResponsive                                int         `json:"isResponsive"`
						       LastVisualChange                            int         `json:"lastVisualChange"`
						       LoadEventEnd                                int         `json:"loadEventEnd"`
						       LoadEventStart                              int         `json:"loadEventStart"`
						       LoadTime                                    int         `json:"loadTime"`
						       MinifySavings                               int         `json:"minify_savings"`
						       MinifyTotal                                 int         `json:"minify_total"`
						       OptimizationChecked                         int         `json:"optimization_checked"`
						       PageSpeedVersion                            int         `json:"pageSpeedVersion"`
						       Render                                      int         `json:"render"`
						       Requests                                    int         `json:"requests"`
						       RequestsDoc                                 int         `json:"requestsDoc"`
						       RequestsFull                                int         `json:"requestsFull"`
						       Responses200                                int         `json:"responses_200"`
						       Responses404                                int         `json:"responses_404"`
						       ResponsesOther                              int         `json:"responses_other"`
						       Result                                      int         `json:"result"`
						       Run                                         int         `json:"run"`
						       ScoreCache                                  int         `json:"score_cache"`
						       ScoreCdn                                    int         `json:"score_cdn"`
						       ScoreCombine                                int         `json:"score_combine"`
						       ScoreCompress                               int         `json:"score_compress"`
						       ScoreCookies                                int         `json:"score_cookies"`
						       ScoreEtags                                  int         `json:"score_etags"`
						       ScoreGzip                                   int         `json:"score_gzip"`
						       ScoreKeep_alive                             int         `json:"score_keep-alive"`
						       ScoreMinify                                 int         `json:"score_minify"`
						       ScoreProgressiveJpeg                        int         `json:"score_progressive_jpeg"`
						       ServerCount                                 int         `json:"server_count"`
						       ServerRtt                                   int         `json:"server_rtt"`
						       TitleTime                                   int         `json:"titleTime"`
						       VisualComplete                              int         `json:"visualComplete"`
					       } `json:"firstView"`
				     RepeatView struct {
						       SpeedIndex                                  int         `json:"SpeedIndex"`
						       TTFB                                        int         `json:"TTFB"`
						       AdultSite                                   int         `json:"adult_site"`
						       Aft                                         int         `json:"aft"`
						       AvgRun                                      interface{} `json:"avgRun"`
						       BytesIn                                     int         `json:"bytesIn"`
						       BytesInDoc                                  int         `json:"bytesInDoc"`
						       BytesOut                                    int         `json:"bytesOut"`
						       BytesOutDoc                                 int         `json:"bytesOutDoc"`
						       Cached                                      int         `json:"cached"`
						       ChromeUserTiming_domComplete                int         `json:"chromeUserTiming.domComplete"`
						       ChromeUserTiming_domContentLoadedEventEnd   int         `json:"chromeUserTiming.domContentLoadedEventEnd"`
						       ChromeUserTiming_domContentLoadedEventStart int         `json:"chromeUserTiming.domContentLoadedEventStart"`
						       ChromeUserTiming_domInteractive             int         `json:"chromeUserTiming.domInteractive"`
						       ChromeUserTiming_domLoading                 int         `json:"chromeUserTiming.domLoading"`
						       ChromeUserTiming_fetchStart                 int         `json:"chromeUserTiming.fetchStart"`
						       ChromeUserTiming_firstContentfulPaint       int         `json:"chromeUserTiming.firstContentfulPaint"`
						       ChromeUserTiming_firstLayout                int         `json:"chromeUserTiming.firstLayout"`
						       ChromeUserTiming_firstPaint                 int         `json:"chromeUserTiming.firstPaint"`
						       ChromeUserTiming_firstTextPaint             int         `json:"chromeUserTiming.firstTextPaint"`
						       ChromeUserTiming_loadEventEnd               int         `json:"chromeUserTiming.loadEventEnd"`
						       ChromeUserTiming_loadEventStart             int         `json:"chromeUserTiming.loadEventStart"`
						       ChromeUserTiming_responseEnd                int         `json:"chromeUserTiming.responseEnd"`
						       ChromeUserTiming_unloadEventEnd             int         `json:"chromeUserTiming.unloadEventEnd"`
						       ChromeUserTiming_unloadEventStart           int         `json:"chromeUserTiming.unloadEventStart"`
						       Connections                                 int         `json:"connections"`
						       Date                                        int         `json:"date"`
						       DocCPUms                                    int         `json:"docCPUms"`
						       DocCPUpct                                   int         `json:"docCPUpct"`
						       DocTime                                     int         `json:"docTime"`
						       DomContentLoadedEventEnd                    int         `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart                  int         `json:"domContentLoadedEventStart"`
						       DomElements                                 int         `json:"domElements"`
						       DomTime                                     int         `json:"domTime"`
						       EffectiveBps                                int         `json:"effectiveBps"`
						       EffectiveBpsDoc                             int         `json:"effectiveBpsDoc"`
						       FirstPaint                                  int         `json:"firstPaint"`
						       FixedViewport                               int         `json:"fixed_viewport"`
						       FullyLoaded                                 int         `json:"fullyLoaded"`
						       FullyLoadedCPUms                            int         `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct                           int         `json:"fullyLoadedCPUpct"`
						       GzipSavings                                 int         `json:"gzip_savings"`
						       GzipTotal                                   int         `json:"gzip_total"`
						       ImageSavings                                int         `json:"image_savings"`
						       ImageTotal                                  int         `json:"image_total"`
						       IsResponsive                                int         `json:"isResponsive"`
						       LastVisualChange                            int         `json:"lastVisualChange"`
						       LoadEventEnd                                int         `json:"loadEventEnd"`
						       LoadEventStart                              int         `json:"loadEventStart"`
						       LoadTime                                    int         `json:"loadTime"`
						       MinifySavings                               int         `json:"minify_savings"`
						       MinifyTotal                                 int         `json:"minify_total"`
						       OptimizationChecked                         int         `json:"optimization_checked"`
						       PageSpeedVersion                            int         `json:"pageSpeedVersion"`
						       Render                                      int         `json:"render"`
						       Requests                                    int         `json:"requests"`
						       RequestsDoc                                 int         `json:"requestsDoc"`
						       RequestsFull                                int         `json:"requestsFull"`
						       Responses200                                int         `json:"responses_200"`
						       Responses404                                int         `json:"responses_404"`
						       ResponsesOther                              int         `json:"responses_other"`
						       Result                                      int         `json:"result"`
						       Run                                         int         `json:"run"`
						       ScoreCache                                  int         `json:"score_cache"`
						       ScoreCdn                                    int         `json:"score_cdn"`
						       ScoreCombine                                int         `json:"score_combine"`
						       ScoreCompress                               int         `json:"score_compress"`
						       ScoreCookies                                int         `json:"score_cookies"`
						       ScoreEtags                                  int         `json:"score_etags"`
						       ScoreGzip                                   int         `json:"score_gzip"`
						       ScoreKeep_alive                             int         `json:"score_keep-alive"`
						       ScoreMinify                                 int         `json:"score_minify"`
						       ScoreProgressiveJpeg                        int         `json:"score_progressive_jpeg"`
						       ServerCount                                 int         `json:"server_count"`
						       ServerRtt                                   int         `json:"server_rtt"`
						       TitleTime                                   int         `json:"titleTime"`
						       VisualComplete                              int         `json:"visualComplete"`
					       } `json:"repeatView"`
			     } `json:"standardDeviation"`
		     SuccessfulFVRuns int    `json:"successfulFVRuns"`
		     SuccessfulRVRuns int    `json:"successfulRVRuns"`
		     Summary          string `json:"summary"`
		     TestURL          string `json:"testUrl"`
		     Tester           string `json:"tester"`
		     TesterDNS        string `json:"testerDNS"`
		     URL              string `json:"url"`
	     } `json:"data"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}


type Response struct {
	Data struct {

		     BwDown       int    `json:"bwDown"`
		     BwUp         int    `json:"bwUp"`
		     Completed    int    `json:"completed"`
		     Connectivity string `json:"connectivity"`
		     From         string `json:"from"`
		     Fvonly       bool   `json:"fvonly"`
		     ID           string `json:"id"`
		     Latency      int    `json:"latency"`
		     Location     string `json:"location"`
		     Median       struct {
				     FirstView struct {
						       TTFB        int    `json:"TTFB"`
						       URL         string `json:"URL"`
						       AdultSite   int    `json:"adult_site"`
						       Aft         int    `json:"aft"`
						       BasePageCdn string `json:"base_page_cdn"`

						       BrowserName      string `json:"browser_name"`
						       BrowserVersion   string `json:"browser_version"`
						       BytesIn          int    `json:"bytesIn"`
						       BytesInDoc       int    `json:"bytesInDoc"`
						       BytesOut         int    `json:"bytesOut"`
						       BytesOutDoc      int    `json:"bytesOutDoc"`
						       Cached           int    `json:"cached"`

						       Connections                                 int `json:"connections"`

						       Date                       int     `json:"date"`
						       DocCPUms                   float64 `json:"docCPUms"`
						       DocCPUpct                  int     `json:"docCPUpct"`
						       DocTime                    int     `json:"docTime"`
						       DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
						       DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
						       DomElements                int     `json:"domElements"`
						       DomTime                    int     `json:"domTime"`

						      FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
						       FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
						       GzipSavings       int     `json:"gzip_savings"`
						       GzipTotal         int     `json:"gzip_total"`
						       ImageSavings      int     `json:"image_savings"`
						       ImageTotal        int     `json:"image_total"`
						       Images            struct {
									   Checklist      string `json:"checklist"`
									   ConnectionView string `json:"connectionView"`
									   ScreenShot     string `json:"screenShot"`
									   Waterfall      string `json:"waterfall"`
								   } `json:"images"`
						       IsResponsive        int    `json:"isResponsive"`
						       LastVisualChange    int    `json:"lastVisualChange"`
						       LoadEventEnd        int    `json:"loadEventEnd"`
						       LoadEventStart      int    `json:"loadEventStart"`
						       LoadTime            int    `json:"loadTime"`
						       MinifySavings       int    `json:"minify_savings"`
						       MinifyTotal         int    `json:"minify_total"`
						       OptimizationChecked int    `json:"optimization_checked"`

						       Title     string `json:"title"`
						       TitleTime int    `json:"titleTime"`
					       } `json:"firstView"`

			     } `json:"median"`


		     Runs   struct {
				     One struct {
						 FirstView struct {
								   TTFB        int    `json:"TTFB"`
								   URL         string `json:"URL"`
								   AdultSite   int    `json:"adult_site"`
								   Aft         int    `json:"aft"`
								   BasePageCdn string `json:"base_page_cdn"`
								   Breakdown   struct {
										       CSS struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"css"`
										       Flash struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"flash"`
										       Font struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"font"`
										       HTML struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"html"`
										       Image struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"image"`
										       Js struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"js"`
										       Other struct {
												   Bytes    int   `json:"bytes"`
												   Color    []int `json:"color"`
												   Requests int   `json:"requests"`
											   } `json:"other"`
									       } `json:"breakdown"`
								   BrowserName      string `json:"browser_name"`
								   BrowserVersion   string `json:"browser_version"`
								   BytesIn          int    `json:"bytesIn"`
								   BytesInDoc       int    `json:"bytesInDoc"`
								   BytesOut         int    `json:"bytesOut"`
								   BytesOutDoc      int    `json:"bytesOutDoc"`
								   Cached           int    `json:"cached"`

								   Connections                                 int `json:"connections"`
								   ConsoleLog                                  []struct {
									   Column             int    `json:"column"`
									   ExecutionContextID int    `json:"executionContextId"`
									   Level              string `json:"level"`
									   Line               int    `json:"line"`
									   Parameters         []struct {
										   Type  string `json:"type"`
										   Value string `json:"value"`
									   } `json:"parameters"`
									   Source string `json:"source"`
									   Stack  struct {
												      CallFrames []struct {
													      ColumnNumber int    `json:"columnNumber"`
													      FunctionName string `json:"functionName"`
													      LineNumber   int    `json:"lineNumber"`
													      ScriptID     string `json:"scriptId"`
													      URL          string `json:"url"`
												      } `json:"callFrames"`
											      } `json:"stack"`
									   Text      string  `json:"text"`
									   Timestamp float64 `json:"timestamp"`
									   Type      string  `json:"type"`
									   URL       string  `json:"url"`
								   } `json:"consoleLog"`
								   Date                       int     `json:"date"`
								   DocCPUms                   float64 `json:"docCPUms"`
								   DocCPUpct                  int     `json:"docCPUpct"`
								   DocTime                    int     `json:"docTime"`
								   DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
								   DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
								   DomElements                int     `json:"domElements"`
								   DomTime                    int     `json:"domTime"`

								   EffectiveBps      int     `json:"effectiveBps"`
								   EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
								   FirstPaint        int     `json:"firstPaint"`
								   FixedViewport     int     `json:"fixed_viewport"`
								   FullyLoaded       int     `json:"fullyLoaded"`
								   FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
								   FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
								   GzipSavings       int     `json:"gzip_savings"`
								   GzipTotal         int     `json:"gzip_total"`
								   ImageSavings      int     `json:"image_savings"`
								   ImageTotal        int     `json:"image_total"`
								   Images            struct {
										       Checklist      string `json:"checklist"`
										       ConnectionView string `json:"connectionView"`
										       ScreenShot     string `json:"screenShot"`
										       Waterfall      string `json:"waterfall"`
									       } `json:"images"`
								   IsResponsive        int    `json:"isResponsive"`
								   LastVisualChange    int    `json:"lastVisualChange"`
								   LoadEventEnd        int    `json:"loadEventEnd"`
								   LoadEventStart      int    `json:"loadEventStart"`
								   LoadTime            int    `json:"loadTime"`
								   MinifySavings       int    `json:"minify_savings"`
								   MinifyTotal         int    `json:"minify_total"`
								   OptimizationChecked int    `json:"optimization_checked"`
								   PageSpeedVersion    string `json:"pageSpeedVersion"`
								   Pages               struct {
										       Breakdown  string `json:"breakdown"`
										       Checklist  string `json:"checklist"`
										       Details    string `json:"details"`
										       Domains    string `json:"domains"`
										       ScreenShot string `json:"screenShot"`
									       } `json:"pages"`
								   RawData struct {
										       Headers      string `json:"headers"`
										       PageData     string `json:"pageData"`
										       RequestsData string `json:"requestsData"`
										       Utilization  string `json:"utilization"`
									       } `json:"rawData"`
								   Render   int `json:"render"`

								   RequestsDoc          int    `json:"requestsDoc"`
								   RequestsFull         int    `json:"requestsFull"`
								   Responses200         int    `json:"responses_200"`
								   Responses404         int    `json:"responses_404"`
								   ResponsesOther       int    `json:"responses_other"`
								   Result               int    `json:"result"`
								   Run                  int    `json:"run"`
								   ScoreCache           int    `json:"score_cache"`
								   ScoreCdn             int    `json:"score_cdn"`
								   ScoreCombine         int    `json:"score_combine"`
								   ScoreCompress        int    `json:"score_compress"`
								   ScoreCookies         int    `json:"score_cookies"`
								   ScoreEtags           int    `json:"score_etags"`
								   ScoreGzip            int    `json:"score_gzip"`
								   ScoreKeep_alive      int    `json:"score_keep-alive"`
								   ScoreMinify          int    `json:"score_minify"`
								   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
								   ServerCount          int    `json:"server_count"`
								   ServerRtt            int    `json:"server_rtt"`
								   Tester               string `json:"tester"`
								   Thumbnails           struct {
										       Checklist  string `json:"checklist"`
										       ScreenShot string `json:"screenShot"`
										       Waterfall  string `json:"waterfall"`
									       } `json:"thumbnails"`
								   Title     string `json:"title"`
								   TitleTime int    `json:"titleTime"`
							   } `json:"firstView"`
						 RepeatView struct {
								   TTFB        int    `json:"TTFB"`
								   URL         string `json:"URL"`
								   AdultSite   int    `json:"adult_site"`
								   Aft         int    `json:"aft"`
								   BasePageCdn string `json:"base_page_cdn"`

								   BrowserName      string `json:"browser_name"`
								   BrowserVersion   string `json:"browser_version"`
								   BytesIn          int    `json:"bytesIn"`
								   BytesInDoc       int    `json:"bytesInDoc"`
								   BytesOut         int    `json:"bytesOut"`
								   BytesOutDoc      int    `json:"bytesOutDoc"`
								   Cached           int    `json:"cached"`

								   Connections                                 int `json:"connections"`
								   ConsoleLog                                  []struct {
									   Column             int    `json:"column"`
									   ExecutionContextID int    `json:"executionContextId"`
									   Level              string `json:"level"`
									   Line               int    `json:"line"`
									   Parameters         []struct {
										   Type  string `json:"type"`
										   Value string `json:"value"`
									   } `json:"parameters"`
									   Source string `json:"source"`
									   Stack  struct {
												      CallFrames []struct {
													      ColumnNumber int    `json:"columnNumber"`
													      FunctionName string `json:"functionName"`
													      LineNumber   int    `json:"lineNumber"`
													      ScriptID     string `json:"scriptId"`
													      URL          string `json:"url"`
												      } `json:"callFrames"`
											      } `json:"stack"`
									   Text      string  `json:"text"`
									   Timestamp float64 `json:"timestamp"`
									   Type      string  `json:"type"`
									   URL       string  `json:"url"`
								   } `json:"consoleLog"`
								   Date                       int     `json:"date"`
								   DocCPUms                   float64 `json:"docCPUms"`
								   DocCPUpct                  int     `json:"docCPUpct"`
								   DocTime                    int     `json:"docTime"`
								   DomContentLoadedEventEnd   int     `json:"domContentLoadedEventEnd"`
								   DomContentLoadedEventStart int     `json:"domContentLoadedEventStart"`
								   DomElements                int     `json:"domElements"`
								   DomTime                    int     `json:"domTime"`

								   EffectiveBps      int     `json:"effectiveBps"`
								   EffectiveBpsDoc   int     `json:"effectiveBpsDoc"`
								   FirstPaint        int     `json:"firstPaint"`
								   FixedViewport     int     `json:"fixed_viewport"`
								   FullyLoaded       int     `json:"fullyLoaded"`
								   FullyLoadedCPUms  float64 `json:"fullyLoadedCPUms"`
								   FullyLoadedCPUpct int     `json:"fullyLoadedCPUpct"`
								   GzipSavings       int     `json:"gzip_savings"`
								   GzipTotal         int     `json:"gzip_total"`
								   ImageSavings      int     `json:"image_savings"`
								   ImageTotal        int     `json:"image_total"`
								   Images            struct {
										       Checklist      string `json:"checklist"`
										       ConnectionView string `json:"connectionView"`
										       ScreenShot     string `json:"screenShot"`
										       Waterfall      string `json:"waterfall"`
									       } `json:"images"`
								   IsResponsive        int    `json:"isResponsive"`
								   LastVisualChange    int    `json:"lastVisualChange"`
								   LoadEventEnd        int    `json:"loadEventEnd"`
								   LoadEventStart      int    `json:"loadEventStart"`
								   LoadTime            int    `json:"loadTime"`
								   MinifySavings       int    `json:"minify_savings"`
								   MinifyTotal         int    `json:"minify_total"`
								   OptimizationChecked int    `json:"optimization_checked"`
								   PageSpeedVersion    string `json:"pageSpeedVersion"`
								   Pages               struct {
										       Breakdown  string `json:"breakdown"`
										       Checklist  string `json:"checklist"`
										       Details    string `json:"details"`
										       Domains    string `json:"domains"`
										       ScreenShot string `json:"screenShot"`
									       } `json:"pages"`
								   RawData struct {
										       Headers      string `json:"headers"`
										       PageData     string `json:"pageData"`
										       RequestsData string `json:"requestsData"`
										       Utilization  string `json:"utilization"`
									       } `json:"rawData"`
								   Domains interface{} `json:"domains"`
								   Render   int `json:"render"`

								   RequestsDoc          int    `json:"requestsDoc"`
								   RequestsFull         int    `json:"requestsFull"`
								   Responses200         int    `json:"responses_200"`
								   Responses404         int    `json:"responses_404"`
								   ResponsesOther       int    `json:"responses_other"`
								   Result               int    `json:"result"`
								   Run                  int    `json:"run"`
								   ScoreCache           int    `json:"score_cache"`
								   ScoreCdn             int    `json:"score_cdn"`
								   ScoreCombine         int    `json:"score_combine"`
								   ScoreCompress        int    `json:"score_compress"`
								   ScoreCookies         int    `json:"score_cookies"`
								   ScoreEtags           int    `json:"score_etags"`
								   ScoreGzip            int    `json:"score_gzip"`
								   ScoreKeep_alive      int    `json:"score_keep-alive"`
								   ScoreMinify          int    `json:"score_minify"`
								   ScoreProgressiveJpeg int    `json:"score_progressive_jpeg"`
								   ServerCount          int    `json:"server_count"`
								   ServerRtt            int    `json:"server_rtt"`
								   Tester               string `json:"tester"`
								   Thumbnails           struct {
										       Checklist  string `json:"checklist"`
										       ScreenShot string `json:"screenShot"`
										       Waterfall  string `json:"waterfall"`
									       } `json:"thumbnails"`
								   Title     string `json:"title"`
								   TitleTime int    `json:"titleTime"`
							   } `json:"repeatView"`
					 } `json:"1"`
			     } `json:"runs"`
		     SuccessfulFVRuns int    `json:"successfulFVRuns"`
		     SuccessfulRVRuns int    `json:"successfulRVRuns"`
		     Summary          string `json:"summary"`
		     TestURL          string `json:"testUrl"`
		     Tester           string `json:"tester"`
		     TesterDNS        string `json:"testerDNS"`
		     URL              string `json:"url"`
	     } `json:"data"`
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
}

type ResultsStorage struct {
	TestID    string      `bson:"testId"`
	Results   string      `bson:"results"`
	LookupUrl string      `bson:"lookupUrl"`
	Url       string      `bson:"url"`
	Tenants   []TenantInfo  `bson:"tenants"`
	CreatedOn time.Time   `bson:"createdOn"`
	UpdatedOn time.Time   `bson:"updatedOn"`
	DeletedOn time.Time   `bson:"deletedOn"`
}

type SystemUrl struct {
	Url       string      `bson:"url"`
	AppId     string      `bson:"appId"`
	DashboardId string    `bson:"dashboardId"`
	CreatedOn time.Time   `bson:"createdOn"`
}

type TestResultsES struct {
	ID		 string             `json:"id"`
	ApplicationId    string		    `json:"applicationid"`
	PageTestUrl      string		    `json:"pagetesturl"`
	TestUrl		 string		    `json:"testurl"`
	Runtime          time.Time	    `json:"runtime"`
	User
        Location         string		    `json:"location"`
	EndpointLocation string		    `json:"endpointlocation"`
	PerformanceScore		    `json:"performancescore"`
	Pageload
	Domains	         []DomainResult	    `json:"domains"`
	ContentTypes     []ContentType      `json:"contenttypes"`
	States		 []State            `json:"states"`
	SuccessfulRuns      int             `json:"successfulruns"`
}

type User struct {
	UserId  string	`json:"userid"`
}

type PerformanceScore struct {
	Cache int	`json:"cache"`
	CDN   int	`json:"cdn"`
	Gzip  int	`json:"gzip"`
	Cookies  int	`json:"cookies"`
	KeepAlive int	`json:"keepalive"`
	Minify   int	`json:"minify"`
	Combine  int	`json:"combine"`
	Compress int	`json:"compress"`
	Etags    int	`json:"etags"`
	Overall  int	`json:"overall"`
}

type Pageload struct {
	Loadtime  int	`json:"loadtime"`
}

type ContentType struct {
	Name     string		`json:"name"`
	Size     int		`json:"size"`
	LoadTime int		`json:"loadtime"`
	Percent  float64	`json:"percent"`
}

type State struct {
	Name        string	`json:"name"`
	TimeSpent   int		`json:"timespent"`
	Percent	    float64     `json:"percent"`
}

type UrlResults struct {
	TenantId  string       `bson:"tenantId"`
	Url    string          `bson:"testUrl"`
	Timestamp []time.Time  `bson:"timestamp"`
	Result string	       `bson:"result"`
}

type TenantInfo struct {
	TenantId  string          `bson:"tenantId"`
	UserId    string          `bson:"userId"`
	Timestamp []time.Time     `bson:"timestamp"`
}

type LocationsResult struct {
	StatusCode int `json:"statusCode"`
	StatusText string `json:"statusText"`
	Data struct {
			   Test struct {
					Label string `json:"Label"`
					Location string `json:"location"`
					Browser string `json:"Browser"`
					RelayServer interface{} `json:"relayServer"`
					RelayLocation interface{} `json:"relayLocation"`
					LabelShort string `json:"labelShort"`
					Default bool `json:"default"`
					PendingTests struct {
						      P1 int `json:"p1"`
						      P2 int `json:"p2"`
						      P3 int `json:"p3"`
						      P4 int `json:"p4"`
						      P5 int `json:"p5"`
						      P6 int `json:"p6"`
						      P7 int `json:"p7"`
						      P8 int `json:"p8"`
						      P9 int `json:"p9"`
						      Total int `json:"Total"`
						      HighPriority int `json:"HighPriority"`
						      LowPriority int `json:"LowPriority"`
						      Testing int `json:"Testing"`
						      Idle int `json:"Idle"`
					      } `json:"PendingTests"`
				} `json:"Test"`
		   } `json:"data"`
}

type LocationsResult1 struct {
	StatusCode int `json:"statusCode"`
	StatusText string `json:"statusText"`
	Data map[string]*json.RawMessage
}

type Tests struct {
	Test struct {
		     Label         string
		     Location      string
		     Browser       string
		     RelayServer   interface{}
		     RelayLocation interface{}
		     LabelShort    string
		     Default       bool
		     PendingTests  struct {
					   P1           int
					   P2           int
					   P3           int
					   P4           int
					   P5           int
					   P6           int
					   P7           int
					   P8           int
					   P9           int
					   Total        int
					   HighPriority int
					   LowPriority  int
					   Testing      int
					   Idle         int
				   }
	     }
	OtherAgents interface{}
}

type SlowestRun struct {
	Loadtime  int		`json:"loadtime"`
	Runtime   int64       	`json:"runtime"`
	Location  string 	`json:"location"`
}

type PageLoadHistory struct {
	LoadtimeHistory   []LoadtimeHistory 	`json:"loadtimeHistory"`
	AverageLoadtime   int   		`json:"averageLoadtime"`
}

type LoadtimeHistory struct {
	Loadtime     int 	`json:"loadtime"`
	Runtime      int64 	`json:"runtime"`
}

type DashboardExists struct {
	Exists		bool	`json:"exists"`
	DashboardId	string  `json:"dashboardId"`
}