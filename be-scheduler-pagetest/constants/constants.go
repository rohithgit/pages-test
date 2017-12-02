// Package constants contains global constants used by other packages in this service.
package constants

// global constants
const (
	SERVICENAME 		string = "PageTest"
	CACHE_PREFIX            string = "spectre-pageload::"
	URL_COUNT            string = "::URLcount"
	PAGELOOKUP_EXPIRATION	int32  = 3600
	TESTRESULTS_EXPIRATION  int32  = 900
	DB_NAME			string = "spectre"
	TEST_COLLECTION 	string = "testresults"
	SPECTRE_TEST_URLS	string = "testurls"
	DB_USERNAME		string = "cloud-user"
	DB_PASSWORD		string = "cloud-user"
	ES_RESULTS_INDEX	string = "pagetestresults"
	ES_RESULTS_TYPE		string = "testresult"
	DEFAULT_LOCATION	string = "Test"
	CODE_403_MESSAGE  = "Forbidden"
	CODE_500_MESSAGE  = "Internal Server Error"
	JSON_CONTENT_TYPE = "application/json"
	SCOPE_WRITE = "write"
	SCOPE_READ = "read"
	SCOPE_SUFFIX = "spectre-be-pagetest:"
	PATH_PREFIX	      		string = "/spectre/v1/pagetest"
	HEALTH_CHECK			string = "/hello"

)
