package game

import "time"

func PuzzleTimeRange() (time.Time, time.Time) {
	// Center the games around Central Time (because that's where I live).
	gmtNow := time.Now().In(time.FixedZone("CT", 0))

	// End with commits from 1.5 years ago so that we likely get authors that people remember.
	startDate := time.Date(gmtNow.Year()-1, gmtNow.Month()-6, gmtNow.Day(), 0, 0, 0, 0, gmtNow.Location())
	// End with commits from a week ago to increase the odds that our user will have an up-to-date history.
	endDate := time.Date(gmtNow.Year(), gmtNow.Month(), gmtNow.Day()-7, 0, 0, 0, 0, gmtNow.Location())

	return startDate, endDate
}
