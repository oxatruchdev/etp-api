package utils

import (
	"strings"
	"time"
)

// FormatTimeInSpanish takes a time and format string, outputs the time in the provided format with Spanish equivalents
func FormatTimeInSpanish(t string, format string) string {
	parsed, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", t)
	if err != nil {
		return ""
	}
	formattedTime := parsed.Format(format)

	// Map of English to Spanish translations for months and days of the week
	months := map[string]string{
		"January":   "Enero",
		"February":  "Febrero",
		"March":     "Marzo",
		"April":     "Abril",
		"May":       "Mayo",
		"June":      "Junio",
		"July":      "Julio",
		"August":    "Agosto",
		"September": "Septiembre",
		"October":   "Octubre",
		"November":  "Noviembre",
		"December":  "Diciembre",
	}

	days := map[string]string{
		"Monday":    "lunes",
		"Tuesday":   "martes",
		"Wednesday": "miércoles",
		"Thursday":  "jueves",
		"Friday":    "viernes",
		"Saturday":  "sábado",
		"Sunday":    "domingo",
	}

	// Replace English month and day names with Spanish equivalents
	for en, es := range months {
		formattedTime = strings.ReplaceAll(formattedTime, en, es)
	}
	for en, es := range days {
		formattedTime = strings.ReplaceAll(formattedTime, en, es)
	}

	return formattedTime
}
