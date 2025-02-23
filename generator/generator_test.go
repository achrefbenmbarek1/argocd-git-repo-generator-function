package generator_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/generator"
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/input/v1beta1"
	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/crossplane-contrib/provider-helm/apis/release/v1beta1"
)

func TestGenerateSetItems(t *testing.T) {
	tests := []struct {
		name string
		urls []string
		want []v1beta1.SetVal
	}{
		{
			name: "Single valid URL",
			urls: []string{"git@github.com:achrefbenmbarek1/libraryManagementBackConfig.git"},
			want: []v1beta1.SetVal{
				{
					Name: "configs.credentialTemplates.github-private-repo-library-management-back-config.sshPrivateKey",
					ValueFrom: &v1beta1.ValueFromSource{
						SecretKeyRef: &v1beta1.DataKeySelector{
							NamespacedName: v1beta1.NamespacedName{
								Namespace: "argocd",
								Name:      "ssh-private-key",
							},
							Key: "ssh-private-key",
						},
					},
				},
			},
		},
		{
			name: "Multiple valid URLs",
			urls: []string{
				"git@github.com:achrefbenmbarek1/repo1.git",
				"git@github.com:achrefbenmbarek1/repo2.git",
			},
			want: []v1beta1.SetVal{
				{
					Name: "configs.credentialTemplates.github-private-repo-repo1.sshPrivateKey",
					ValueFrom: &v1beta1.ValueFromSource{
						SecretKeyRef: &v1beta1.DataKeySelector{
							NamespacedName: v1beta1.NamespacedName{
								Namespace: "argocd",
								Name:      "ssh-private-key",
							},
							Key: "ssh-private-key",
						},
					},
				},
				{
					Name: "configs.credentialTemplates.github-private-repo-repo2.sshPrivateKey",
					ValueFrom: &v1beta1.ValueFromSource{
						SecretKeyRef: &v1beta1.DataKeySelector{
							NamespacedName: v1beta1.NamespacedName{
								Namespace: "argocd",
								Name:      "ssh-private-key",
							},
							Key: "ssh-private-key",
						},
					},
				},
			},
		},
		{
			name: "Invalid URL skipped",
			urls: []string{"invalid-url"},
			want: []v1beta1.SetVal{},
		},
		{
			name: "Mix of valid and invalid URLs",
			urls: []string{
				"git@github.com:achrefbenmbarek1/valid-repo.git",
				"invalid-url",
				"git@github.com:achrefbenmbarek1/another-valid.git",
			},
			want: []v1beta1.SetVal{
				{
					Name: "configs.credentialTemplates.github-private-repo-valid-repo.sshPrivateKey",
					ValueFrom: &v1beta1.ValueFromSource{
						SecretKeyRef: &v1beta1.DataKeySelector{
							NamespacedName: v1beta1.NamespacedName{
								Namespace: "argocd",
								Name:      "ssh-private-key",
							},
							Key: "ssh-private-key",
						},
					},
				},
				{
					Name: "configs.credentialTemplates.github-private-repo-another-valid.sshPrivateKey",
					ValueFrom: &v1beta1.ValueFromSource{
						SecretKeyRef: &v1beta1.DataKeySelector{
							NamespacedName: v1beta1.NamespacedName{
								Namespace: "argocd",
								Name:      "ssh-private-key",
							},
							Key: "ssh-private-key",
						},
					},
				},
			},
		},
		{
			name: "Empty input",
			urls: []string{},
			want: []v1beta1.SetVal{},
		},
	}

	// Ignore unexported fields in comparison
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(v1beta1.SetVal{}),
		cmpopts.IgnoreUnexported(v1beta1.ValueFromSource{}),
		cmpopts.IgnoreUnexported(v1beta1.DataKeySelector{}),
		cmpopts.IgnoreUnexported(v1beta1.NamespacedName{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generator.GenerateSetItems(tt.urls)

			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("GenerateSetItems() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateValues(t *testing.T) {
	in := &v1beta1input.ArgoRepoGeneratorInput{
		Spec: v1beta1input.ArgoRepoGeneratorSpec{
			ForProvider: v1beta1input.ForProviderSpec{
				Values: v1beta1input.ValuesSpec{
					Server: v1beta1input.ServerSpec{
						Service: v1beta1input.ServerServiceSpec{
							Type: "ClusterIP",
						},
						Ingress: v1beta1input.IngressSpec{
							Enabled:          true,
							IngressClassName: "nginx",
							Hostname:         "argocd.example.com",
							Path:             "/",
							PathType:         "Prefix",
							TLS:              true,
						},
					},
					Configs: v1beta1input.ConfigSpec{
						Secret: v1beta1input.SecretSpec{
							CreateSecret:              true,
							ArgoCDServerAdminPassword: "admin123",
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name      string
		in        *v1beta1input.ArgoRepoGeneratorInput
		reposUrls []string
	}{
		{
			name:      "Single valid URL",
			in:        in,
			reposUrls: []string{"git@github.com:user/repo1.git"},
		},
		{
			name: "Multiple valid URLs",
			in:   in,
			reposUrls: []string{
				"git@github.com:user/repo1.git",
				"git@github.com:user/repo2.git",
			},
		},
		{
			name:      "Invalid URL skipped",
			in:        in,
			reposUrls: []string{"invalid-url"},
		},
		{
			name: "Mix of valid and invalid URLs",
			in:   in,
			reposUrls: []string{
				"git@github.com:user/valid-repo.git",
				"invalid-url",
				"git@github.com:user/another-valid.git",
			},
		},
		{
			name:      "Empty input",
			in:        in,
			reposUrls: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := generateExpectedValues(tt.in, tt.reposUrls)

			gotBytes, err := generator.GenerateValues(tt.in, tt.reposUrls)
			if err != nil {
				t.Fatalf("GenerateValues() returned error: %v", err)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(gotBytes, &got); err != nil {
				t.Fatalf("Failed to unmarshal got JSON: %v", err)
			}

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("GenerateValues() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func generateExpectedValues(in *v1beta1input.ArgoRepoGeneratorInput, reposUrls []string) map[string]interface{} {
	credentialTemplates := make(map[string]interface{})
	repositories := make(map[string]interface{})

	for _, url := range reposUrls {
		if !isValidURL(url) {
			continue
		}

		repoName, username := parseRepoInfo(url)
		if repoName == "" || username == "" {
			continue
		}

		key := "github-private-repo-" + utils.ToKebabCaseForValidKuberSecrets(repoName)
		credentialTemplates[key] = map[string]interface{}{"url": url}
		repositories[key] = map[string]interface{}{
			"url":            url,
			"type":           "git",
			"name":           utils.ToKebabCaseForValidKuberSecrets(repoName),
			"credentialName": username,
		}
	}

	// Handle nil slice conversion for extraArgs
	var extraArgs interface{}
	if in.Spec.ForProvider.Values.Server.ExtraArgs != nil {
		extraArgs = in.Spec.ForProvider.Values.Server.ExtraArgs
	}

	return map[string]interface{}{
		"server": map[string]interface{}{
			"service": map[string]interface{}{
				"type": in.Spec.ForProvider.Values.Server.Service.Type,
			},
			"ingress": map[string]interface{}{
				"enabled":          in.Spec.ForProvider.Values.Server.Ingress.Enabled,
				"ingressClassName": in.Spec.ForProvider.Values.Server.Ingress.IngressClassName,
				"hostname":         in.Spec.ForProvider.Values.Server.Ingress.Hostname,
				"path":             in.Spec.ForProvider.Values.Server.Ingress.Path,
				"pathType":         in.Spec.ForProvider.Values.Server.Ingress.PathType,
				"tls":              in.Spec.ForProvider.Values.Server.Ingress.TLS,
			},
			"extraArgs": extraArgs, // Now matches JSON null when empty
		},
		"configs": map[string]interface{}{
			"secret": map[string]interface{}{
				"createSecret":              in.Spec.ForProvider.Values.Configs.Secret.CreateSecret,
				"argocdServerAdminPassword": in.Spec.ForProvider.Values.Configs.Secret.ArgoCDServerAdminPassword,
			},
			"credentialTemplates": credentialTemplates,
			"repositories":        repositories,
		},
	}
}

// isValidURL checks if the URL matches the expected format (simplified for testing)
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "git@github.com:") && strings.HasSuffix(url, ".git")
}

// parseRepoInfo extracts repo name and username from a valid URL (simplified for testing)
func parseRepoInfo(url string) (repoName, username string) {
	if !isValidURL(url) {
		return "", ""
	}

	// Trim prefix and suffix
	trimmed := strings.TrimPrefix(url, "git@github.com:")
	trimmed = strings.TrimSuffix(trimmed, ".git")

	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		return "", ""
	}

	username = parts[0]
	repoName = parts[1]
	return repoName, username
}
