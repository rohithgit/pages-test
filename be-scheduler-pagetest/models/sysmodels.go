package models

type AuthStruct struct {
	Type        string `json:"-"`
	User        string `json:"-"`
	Password    string `json:"-"`
	AuthToken   string `json:"-"`
	AuthHeader  string `json:"-"`
	ContentType string `json:"-"`
	AccessToken string `json:"-"`
}

type RequestHeader struct {
	Key	string
	Value	string
}

type ErrorCode struct {
	Code  string `json:"errorCode"`
	Error string `json:"errorString"`
}