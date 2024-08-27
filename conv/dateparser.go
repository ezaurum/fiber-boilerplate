package conv

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

var (
	// 2006/1/2, 2006-01-02, 2006.01.02, 2006.1.2, 2006년 1월 2일
	yearDate = regexp.MustCompile("([0-9]{4})\\s*[/\\-.년 ]\\s*([0-9]{1,2})\\s*[/\\-.월 ]\\s*([0-9]{1,2})")

	// 15:04
	timeReg = regexp.MustCompile("\\s*([0-2]?[0-9])\\s*[:시 ]\\s*([0-6]?[0-9])")

	// 2006/1/2, 2006-01-02, 2006.01.02, 2006.1.2, 2006년 1월 2일
	noYearDate = regexp.MustCompile("^([0-9]{1,2})\\s*[/\\-.월 ]\\s*([0-9]{1,2})")
)

// ParseHumanDateTime 대충 만든 문자열을 날짜/시간으로 파싱
func ParseHumanDateTime(numberString string) (*time.Time, error) {
	space := strings.TrimSpace(numberString)
	if !timeReg.MatchString(space) {
		return nil, errors.New(fmt.Sprintf("cannot parse time %v as 15:04", numberString))
	}
	timeStringMatched := timeReg.FindStringSubmatch(space)

	return parseTime(space, timeStringMatched)
}

// ParseHumanDate 대충 만든 문자열을 날짜로 파싱
func ParseHumanDate(numberString string) (*time.Time, error) {
	space := strings.TrimSpace(numberString)
	return parseTime(space, []string{"0", "0", "0"})
}

// 날짜 파싱
func parseTime(fullString string, timeStringMatched []string) (*time.Time, error) {
	var dateStringMatched []string
	switch {
	case yearDate.MatchString(fullString):
		dateStringMatched = yearDate.FindStringSubmatch(fullString)
	case noYearDate.MatchString(fullString):
		dateStringMatched = append([]string{"0", fmt.Sprint(time.Now().Year())}, noYearDate.FindStringSubmatch(fullString)[1:]...)
	default:
		return nil, errors.New(fmt.Sprintf("cannot parse %v", fullString))
	}

	result := fmt.Sprintf("%04v%02v%02vT%02v%02v",
		dateStringMatched[1], dateStringMatched[2], dateStringMatched[3],
		timeStringMatched[1], timeStringMatched[2])

	if t, e := time.ParseInLocation("20060102T1504", result, time.Now().Location()); nil == e {
		return &t, nil
	} else {
		return nil, errors.Wrapf(e, "cannot parse %v", fullString)
	}
}
