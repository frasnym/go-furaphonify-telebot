package common

import "regexp"

// RemoveNonNumeric takes a string and removes all non-numeric characters (anything that is not a digit).
// It uses a regular expression to perform the removal and returns the modified string with only numeric characters.
func RemoveNonNumeric(str string) string {
	regex := regexp.MustCompile(`[^0-9]+`)
	return regex.ReplaceAllString(str, "")
}
