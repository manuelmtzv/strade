package utils

import (
	"strings"
	"time"
)

func ParseSpanishDate(dateStr string) (time.Time, error) {
	monthMap := map[string]string{
		"enero":      "January",
		"febrero":    "February",
		"marzo":      "March",
		"abril":      "April",
		"mayo":       "May",
		"junio":      "June",
		"julio":      "July",
		"agosto":     "August",
		"septiembre": "September",
		"octubre":    "October",
		"noviembre":  "November",
		"diciembre":  "December",
	}

	for spanish, english := range monthMap {
		dateStr = strings.ReplaceAll(strings.ToLower(dateStr), spanish, english)
	}

	dateStr = strings.ReplaceAll(dateStr, " de ", " ")

	return time.Parse("January 2 2006", dateStr)
}
