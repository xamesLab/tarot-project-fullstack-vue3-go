package utils

import "time"

func ValidateTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

func GetAccountingDate(timezone string) (string, error) {
	return GetAccountingDateAt(time.Now(), timezone)
}

func GetAccountingDateAt(now time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	return getAccountingDateForTime(now.In(loc)), nil
}

func getAccountingDateForTime(now time.Time) string {
	if now.Hour() < 3 {
		now = now.Add(-24 * time.Hour)
	}

	return now.Format("2006-01-02")
}

func GetNextResetTime(timezone string) (time.Time, error) {
	return GetNextResetTimeAt(time.Now(), timezone)
}

func GetNextResetTimeAt(now time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return getNextResetTimeForTime(now.In(loc)), nil
}

func getNextResetTimeForTime(now time.Time) time.Time {
	reset := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	if !now.Before(reset) {
		reset = reset.Add(24 * time.Hour)
	}

	return reset
}
