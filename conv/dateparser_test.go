package conv

import (
	"testing"
	"time"
)

func TestParseHumanDate_1월2일(t *testing.T) {
	patterns := []string{
		"1월 2일",
		"1월2일",
		"1 / 2",
		"1월2일",
		"1.2",
		"1  .  2  .",
		"1   2 ",
	}
	for _, p := range patterns {
		date, err := ParseHumanDate(p)
		if err != nil {
			t.Error(err)
		}
		if date == nil {
			t.Error("date is nil", p)
		}
		year := time.Now().Year()
		if date.Year() != year {
			t.Error("date is not this year", date.Year(), p)
		}
		if date.Format("01-02") != "01-02" {
			t.Error("date is not 01-02", date.Format("2006-01-02"), p)
		}
	}
}

func TestParseHumanDate_20241월2일(t *testing.T) {
	patterns := []string{
		" 2024 년 1월 2일",
		"2024.1월2일",
		"2024 1 / 2",
		"2024-1월2일",
		"2024.1.2",
		"2024   1  .  2  .",
	}
	for _, p := range patterns {
		date, err := ParseHumanDate(p)
		if err != nil {
			t.Error(err)
		}
		if date == nil {
			t.Error("date is nil", p)
		}
		year := time.Now().Year()
		if date.Year() != year {
			t.Error("date is not this year", date.Year(), p)
		}
		if date.Format("01-02") != "01-02" {
			t.Error("date is not 01-02", date.Format("2006-01-02"), p)
		}
	}
}
