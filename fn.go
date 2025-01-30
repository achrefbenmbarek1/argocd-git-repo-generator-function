package main

import (
	"context"

	"github.com/achrefbenmbarek1/argocd-git-repo-generator-function/generator"
	v1beta1input "github.com/achrefbenmbarek1/argocd-git-repo-generator-function/input/v1beta1"

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

type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

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
