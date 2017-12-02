package helper

import (
	"testing"
	"fmt"
)

func Test1GradeF(t *testing.T) {
	percentStr := "-1"
	grade, err := getGrade(percentStr)
	fmt.Printf("Grade is %s\n", grade)
	fmt.Printf("Error is %v\n", err)
	if grade != "F" {
		fmt.Printf("Percent is %s ", percentStr)
		t.Error("Expected grade F")
	}
}

func Test2GradeF(t *testing.T) {
	percentStr := "0"
	grade, err := getGrade(percentStr)
	fmt.Printf("Grade is %s\n", grade)
	fmt.Printf("Error is %v\n", err)
	if grade != "F" {
		fmt.Printf("Percent is %s ", percentStr)
		t.Error("Expected grade F")
	}
}

func Test3GradeF(t *testing.T) {
	percentStr := "59"
	grade, err := getGrade(percentStr)
	fmt.Printf("Grade is %s\n", grade)
	fmt.Printf("Error is %v\n", err)
	if grade != "F" {
		fmt.Printf("Percent is %s ", percentStr)
		t.Error("Expected grade F")
	}
}

func Test1GradeD(t *testing.T) {
	percentStr := "60"
	grade, err := getGrade(percentStr)
	fmt.Printf("Grade is %s\n", grade)
	fmt.Printf("Error is %v\n", err)
	if grade != "D" {
		fmt.Printf("Percent is %s ", percentStr)
		t.Error("Expected grade D")
	}
}

func Test2GradeD(t *testing.T) {
	percentStr := "69"
	grade, err := getGrade(percentStr)
	fmt.Printf("Grade is %s\n", grade)
	fmt.Printf("Error is %v\n", err)
	if grade != "D" {
		fmt.Printf("Percent is %s ", percentStr)
		t.Error("Expected grade D")
	}
}

