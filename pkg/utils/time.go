package utils

import "time"

// IsValidDate checks if a date string is in the format YYYY-MM-DD
func IsValidDate(date string) bool {
    _, err := time.Parse("2006-01-02", date)
    return err == nil
}