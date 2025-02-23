package generator

import (
	"encoding/json"
	"fmt"

	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/input/v1beta1"
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/utils"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/crossplane-contrib/provider-helm/apis/release/v1beta1"
)

func GenerateValues(in *v1beta1input.ArgoRepoGeneratorInput, reposUrls []string) ([]byte, error) {
	values := map[string]interface{}{
		"server": in.Spec.ForProvider.Values.Server,
		"configs": map[string]interface{}{
			"secret":              in.Spec.ForProvider.Values.Configs.Secret,
			"credentialTemplates": GenerateCredentialTemplates(reposUrls),
			"repositories":        GenerateRepositories(reposUrls),
		},
	}
	return json.Marshal(values)
}

func GenerateSetItems(urls []string) []v1beta1.SetVal {
	setItems := []v1beta1.SetVal{}
	for _, url := range urls {
		repoName, _, valid := utils.ExtractRepoInfo(url)
		if !valid {
			continue
		}
		validResourceNameForKuber := utils.ToKebabCaseForValidKuberSecrets(repoName)
		key := "github-private-repo-" + validResourceNameForKuber
		setItems = append(setItems, v1beta1.SetVal{
			Name: fmt.Sprintf("configs.credentialTemplates.%s.sshPrivateKey", key),
			ValueFrom: &v1beta1.ValueFromSource{
				SecretKeyRef: &v1beta1.DataKeySelector{
					NamespacedName: v1beta1.NamespacedName{
						Namespace: "argocd",
						Name:      "ssh-private-key",
					},
					Key: "ssh-private-key",
				},
			},
		})
	}
	return setItems
}

func GenerateHelmRelease(in *v1beta1input.ArgoRepoGeneratorInput, reposUrls []string, valuesJSON []byte) (*v1beta1.Release, error) {
	helmRelease := &v1beta1.Release{
		Spec: v1beta1.ReleaseSpec{
			ForProvider: v1beta1.ReleaseParameters{
				Chart: v1beta1.ChartSpec{
					Name:       in.Spec.ForProvider.Chart.Name,
					Repository: in.Spec.ForProvider.Chart.Repository,
					Version:    in.Spec.ForProvider.Chart.Version,
				},
				Namespace: in.Spec.ForProvider.Namespace,
				ValuesSpec: v1beta1.ValuesSpec{
					Values: runtime.RawExtension{Raw: valuesJSON},
					Set:    GenerateSetItems(reposUrls),
				},
			},
		},
	}

	helmRelease.SetProviderConfigReference(&v1.Reference{Name: "helm-provider"})
	return helmRelease, nil
}
