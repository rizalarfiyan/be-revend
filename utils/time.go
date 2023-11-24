package utils

import "time"

func RemaniningToday() time.Duration {
	now := time.Now()
	dateTomorrow := now.AddDate(0, 0, 1)
	tomorrow := time.Date(dateTomorrow.Year(), dateTomorrow.Month(), dateTomorrow.Day(), 0, 0, 0, 0, dateTomorrow.Location())
	return tomorrow.Sub(now)
}
