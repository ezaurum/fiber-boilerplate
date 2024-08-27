package conv

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

var (
	// 2006/1/2
	slashDate = regexp.MustCompile("([0-9]{4})/([0-9]{1,2})/([0-9]{1,2})")

	// 2006-01-02
	fullISO = regexp.MustCompile("([0-9]{4})\\s*-([0-1][0-9])\\s*-([0-3][0-9])")
	// 2006-1-2
	shortISO = regexp.MustCompile("([0-9]{4})\\s*-(1?[0-9])\\s*-([1-3]?[0-9])")

	// 2006.01.02
	fullDot = regexp.MustCompile("([0-9]{4})\\s*.\\s*([0-1][0-9])\\s*.\\s*([0-3][0-9])")
	// 2006.1.2
	shortDot = regexp.MustCompile("([0-9]{4})\\s*.\\s*([1]?[0-9])\\s*.\\s*([1-3]?[0-9])")

	// 15:04
	timeReg = regexp.MustCompile("\\s*([0-2]?[0-9])\\s*:\\s*([0-6]?[0-9])")
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
	case slashDate.MatchString(fullString):
		dateStringMatched = slashDate.FindStringSubmatch(fullString)
	case fullISO.MatchString(fullString):
		dateStringMatched = fullISO.FindStringSubmatch(fullString)
	case shortISO.MatchString(fullString):
		dateStringMatched = shortISO.FindStringSubmatch(fullString)
	case fullDot.MatchString(fullString):
		dateStringMatched = fullDot.FindStringSubmatch(fullString)
	case shortDot.MatchString(fullString):
		dateStringMatched = shortDot.FindStringSubmatch(fullString)
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
