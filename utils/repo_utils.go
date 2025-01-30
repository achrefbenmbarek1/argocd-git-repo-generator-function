package utils

import "regexp"

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
