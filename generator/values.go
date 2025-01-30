package generator

import (
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/models"
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/utils"
)

func GenerateRepositories(urls []string) map[string]interface{} {
	repositories := make(map[string]interface{})
	for _, url := range urls {
		repoName, username, valid := utils.ExtractRepoInfo(url)
		if !valid {
			continue
		}
		key := "github-private-repo-" + repoName
		repositories[key] = models.Repository{
			URL:            url,
			Type:           "git",
			Name:           repoName,
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
		key := "github-private-repo-" + repoName
		credentials[key] = map[string]interface{}{
			"url": url,
		}
	}
	return credentials
}

