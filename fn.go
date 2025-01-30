package main

import (
	"context"

	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/generator"
	v1beta1input "github.com/achrefbenmbarek1/argocd-git-repo-generator-function/input/v1beta1Input"

	// "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"

	"github.com/crossplane-contrib/provider-helm/apis/release/v1beta1"
)

// func generateHelmRelease(in *v1beta1Input.ArgoRepoGeneratorInput, reposUrls []string, valuesJSON []byte) (*v1beta1.Release, error) {
// 	helRelease := &v1beta1.Release{
// 		Spec: v1beta1.ReleaseSpec{
// 			ForProvider: v1beta1.ReleaseParameters{
// 				Chart: v1beta1.ChartSpec{
// 					Name:       in.Spec.ForProvider.Chart.Name,
// 					Repository: in.Spec.ForProvider.Chart.Repository,
// 					Version:    in.Spec.ForProvider.Chart.Version,
// 				},
//
// 				Namespace: in.Spec.ForProvider.Namespace,
// 				ValuesSpec: v1beta1.ValuesSpec{
// 					Values: runtime.RawExtension{Raw: valuesJSON},
// 					Set:    generateSetItems(reposUrls),
// 				},
// 			},
// 		},
// 	}
//
// 	helRelease.SetProviderConfigReference(&v1.Reference{Name: "helm-provider"})
// 	return helRelease, nil
// }
//
// func generateValues(in *v1beta1Input.ArgoRepoGeneratorInput, reposUrls []string) ([]byte, error) {
// 	values := map[string]interface{}{
// 		"server": in.Spec.ForProvider.Values.Server,
// 		"configs": map[string]interface{}{
// 			"secret":              in.Spec.ForProvider.Values.Configs.Secret,
// 			"credentialTemplates": generateCredentialTemplates(reposUrls),
// 			"repositories":        generateRepositories(reposUrls),
// 		},
// 	}
// 	return json.Marshal(values)
// }
//
// func extractRepoInfo(url string) (string, string, bool) {
// 	repoRegex := regexp.MustCompile(`.*/([^/]+)\.git$`)
// 	userRegex := regexp.MustCompile(`github\.com:([^/]+)`)
//
// 	repoMatches := repoRegex.FindStringSubmatch(url)
// 	userMatches := userRegex.FindStringSubmatch(url)
//
// 	if len(repoMatches) < 2 || len(userMatches) < 2 {
// 		return "", "", false
// 	}
// 	return repoMatches[1], userMatches[1], true
// }
//
// func generateRepositories(urls []string) map[string]interface{} {
// 	repositories := make(map[string]interface{})
// 	for _, url := range urls {
// 		repoName, username, valid := extractRepoInfo(url)
// 		if !valid {
// 			continue
// 		}
// 		key := "github-private-repo-" + repoName
// 		repositories[key] = Repository{
// 			URL:            url,
// 			Type:           "git",
// 			Name:           repoName,
// 			CredentialName: username,
// 		}
// 	}
// 	return repositories
// }
//
// func generateCredentialTemplates(urls []string) map[string]interface{} {
// 	credentials := make(map[string]interface{})
// 	for _, url := range urls {
// 		repoName, _, valid := extractRepoInfo(url)
// 		if !valid {
// 			continue
// 		}
// 		key := "github-private-repo-" + repoName
// 		credentials[key] = map[string]interface{}{
// 			"url": url,
// 		}
// 	}
// 	return credentials
// }
//
// func generateSetItems(urls []string) []v1beta1.SetVal {
// 	var setItems []v1beta1.SetVal
// 	for _, url := range urls {
// 		repoName, _, valid := extractRepoInfo(url)
// 		if !valid {
// 			continue
// 		}
// 		key := "github-private-repo-" + repoName
// 		setItems = append(setItems, v1beta1.SetVal{
// 			Name: fmt.Sprintf("configs.credentialTemplates.%s.sshPrivateKey", key),
// 			ValueFrom: &v1beta1.ValueFromSource{
// 				SecretKeyRef: &v1beta1.DataKeySelector{
// 					NamespacedName: v1beta1.NamespacedName{
// 						Namespace: "argocd",
// 						Name:      "ssh-private-key",
// 					},
// 					Key: "ssh-private-key",
// 				},
// 			},
// 		})
// 	}
// 	return setItems
// }

type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// type Repository struct {
// 	URL            string `json:"url" yaml:"url"`
// 	Type           string `json:"type" yaml:"type"`
// 	Name           string `json:"name" yaml:"name"`
// 	CredentialName string `json:"credentialName" yaml:"credentialName"`
// }

func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {

	f.log.Info("Running function", "tag", req.GetMeta().GetTag())
	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1input.ArgoRepoGeneratorInput{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}

	xr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource from %T", req))
		return rsp, nil
	}

	reposUrls, err := xr.Resource.GetStringArray("spec.repos-urls")
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get repos-urls from %T", xr))
		return rsp, nil
	}

	valuesJSON, _ := generator.GenerateValues(in, reposUrls)

	_ = v1beta1.SchemeBuilder.AddToScheme(composed.Scheme)

	helmRelease, _ := generator.GenerateHelmRelease(in, reposUrls, valuesJSON)

	cd, err := composed.From(helmRelease)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot convert %T to %T", helmRelease, &composed.Unstructured{}))
		return rsp, nil
	}

	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired resources from %T", req))
		return rsp, nil
	}
	desired[resource.Name(in.Name)] = &resource.DesiredComposed{Resource: cd}

	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
		return rsp, nil
	}

	response.ConditionTrue(rsp, "FunctionSuccess", "Success").TargetCompositeAndClaim()
	return rsp, nil
}
