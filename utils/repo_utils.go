package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func ExtractRepoInfo(url string) (string, string, bool) {
	repoRegex := regexp.MustCompile(`.*/([^/]+)\.git$`)
	userRegex := regexp.MustCompile(`github\.com:([^/]+)`)

	repoMatches := repoRegex.FindStringSubmatch(url)
	userMatches := userRegex.FindStringSubmatch(url)

	if len(repoMatches) < 2 || len(userMatches) < 2 {
		return "", "", false
	}
	return repoMatches[1], userMatches[1], true
}

func ToKebabCaseForValidKuberSecrets(input string) string {
	var result strings.Builder

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-') // Add '-' before uppercase letters
			}
			result.WriteRune(unicode.ToLower(r)) // Convert uppercase to lowercase
		} else {
			result.WriteRune(r) // Keep lowercase letters unchanged
		}
	}

	return result.String()
}
