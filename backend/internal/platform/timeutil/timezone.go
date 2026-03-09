package timeutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DefaultTimezone = "UTC+00:00"

var utcOffsetPattern = regexp.MustCompile(`^UTC([+-])(\d{2}):(\d{2})$`)

func NormalizeTimezone(input string) (string, error) {
	timezone := strings.TrimSpace(input)
	if timezone == "" || strings.EqualFold(timezone, "utc") {
		return DefaultTimezone, nil
	}

	if _, err := parseUTCOffsetSeconds(timezone); err == nil {
		return timezone, nil
	}

	if _, err := time.LoadLocation(timezone); err == nil {
		return timezone, nil
	}

	return "", errors.New("invalid timezone")
}

func LocationFromTimezone(timezone string) (*time.Location, error) {
	tz, err := NormalizeTimezone(timezone)
	if err != nil {
		return nil, err
	}

	if seconds, err := parseUTCOffsetSeconds(tz); err == nil {
		return time.FixedZone(tz, seconds), nil
	}
	return time.LoadLocation(tz)
}

func BusinessDateAtReset(now time.Time, timezone string, resetHour int) (string, error) {
	loc, err := LocationFromTimezone(timezone)
	if err != nil {
		return "", err
	}

	localNow := now.In(loc)
	dayAnchor := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), resetHour, 0, 0, 0, loc)
	businessDate := localNow
	if localNow.Before(dayAnchor) {
		businessDate = businessDate.AddDate(0, 0, -1)
	}

	return businessDate.Format("2006-01-02"), nil
}

func parseUTCOffsetSeconds(tz string) (int, error) {
	matches := utcOffsetPattern.FindStringSubmatch(tz)
	if len(matches) != 4 {
		return 0, errors.New("not a UTC offset timezone")
	}

	hours, _ := strconv.Atoi(matches[2])
	minutes, _ := strconv.Atoi(matches[3])
	if hours > 14 || minutes > 59 || (hours == 14 && minutes != 0) {
		return 0, fmt.Errorf("timezone out of range: %s", tz)
	}

	sign := 1
	if matches[1] == "-" {
		sign = -1
	}

	return sign * ((hours * 3600) + (minutes * 60)), nil
}
