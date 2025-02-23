package generator

import (
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/models"
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/utils"
)

// GenerateRepositories creates a map of repository configurations from Git URLs
// Parses SSH-style URLs (git@github.com:user/repo.git) to extract repo metadata
func GenerateRepositories(urls []string) map[string]interface{} {
	repositories := make(map[string]interface{})
	for _, url := range urls {
		repoName, username, valid := utils.ExtractRepoInfo(url)
		if !valid {
			continue
		}
		validResourceNameForKuber := utils.ToKebabCaseForValidKuberSecrets(repoName)
		key := "github-private-repo-" + validResourceNameForKuber
		repositories[key] = models.Repository{
			URL:            url,
			Type:           "git",
			Name:           validResourceNameForKuber,
			CredentialName: username,
		}
	}
	return repositories
}

func GenerateCredentialTemplates(urls []string) map[string]interface{} {
	credentials := make(map[string]interface{})
	for _, url := range urls {
		repoName, _, valid := utils.ExtractRepoInfo(url)
		if !valid {
			continue
		}
		validResourceNameForKuber := utils.ToKebabCaseForValidKuberSecrets(repoName)
		key := "github-private-repo-" + validResourceNameForKuber
		credentials[key] = map[string]interface{}{
			"url": url,
		}
	}
	return credentials
}
