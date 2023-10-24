package common

import "regexp"

// RemoveNonAlphanumeric takes a string and removes all non-alphanumeric characters (anything that is not a letter or a digit).
// It uses a regular expression to perform the removal and returns the modified string with only alphanumeric characters.
func RemoveNonAlphanumeric(str string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return regex.ReplaceAllString(str, "")
}
