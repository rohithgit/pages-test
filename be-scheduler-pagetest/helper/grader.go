package helper

import (

	"errors"
)

//Letter Grade Percent
//A: 	90%+
//B: 	80–89%
//C: 	70–79%
//D:	60–69%
//F:	0–59%

func GetGrade(percent int) (string, error) {

	if percent >= 90  {
		return "A", nil
	}
	if percent >= 80 && percent <= 89 {
		return "B", nil
	}
	if percent >= 70 && percent <= 79 {
		return "C", nil
	}
	if percent >= 60 && percent <= 69 {
		return "D", nil
	}

	if (percent >= 0  && percent <= 59)  {
		return "F", nil
	}

	if percent == -1 {
		return "N", nil
	}
	return "", errors.New("Cannot calculate grade")
}
