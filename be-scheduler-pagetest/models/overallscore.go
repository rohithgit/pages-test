package models

type ScoreDetailResult struct {
	Overallscore string `json:"overallScore"`
	Individualscore []ScoreDetail  `json:"individualScore"`
}
