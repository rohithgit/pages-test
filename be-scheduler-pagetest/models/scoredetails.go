package models

type ScoreDetail struct {
	Category    string `json:"category"`
	Score int `json:"score"`
	Suggestion string `json:"suggestion"`
	SuggestionDetails string `json:"suggestionDetails"`
}

type ScoreHistory struct {
	Score     int 	`json:"score"`
	Runtime      int64 	`json:"runtime"`
}
