package main
//
// import (
// 	"context"
// 	"testing"
// 	"time"
//
// 	"github.com/google/go-cmp/cmp"
// 	"github.com/google/go-cmp/cmp/cmpopts"
// 	"google.golang.org/protobuf/testing/protocmp"
// 	"google.golang.org/protobuf/types/known/durationpb"
//
// 	"github.com/crossplane/crossplane-runtime/pkg/logging"
// 	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
// 	"github.com/crossplane/function-sdk-go/resource"
// )
//
// func TestRunFunction(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *fnv1.RunFunctionRequest
// 	}
// 	type want struct {
// 		rsp *fnv1.RunFunctionResponse
// 		err error
// 	}
//
// 	cases := map[string]struct {
// 		reason string
// 		args   args
// 		want   want
// 	}{
// 		"CreateArgoCDHelmRelease": {
// 			reason: "The Function should create a Helm Release for Argo CD with repository configurations",
// 			args: args{
// 				req: &fnv1.RunFunctionRequest{
// 					Observed: &fnv1.State{
// 						Composite: &fnv1.Resource{
// 							Resource: resource.MustStructJSON(`{
// 								"apiVersion": "example.crossplane.io/v1",
// 								"kind": "XR",
// 								"metadata": {
// 									"name": "example-xr"
// 								},
// 								"spec": {
// 									"reposUrls": [
// 										"git@github.com:achrefbenmbarek1/libraryManagementBackConfig.git",
// 										"git@github.com:achrefbenmbarek1/libraryManagementAppOfApps.git"
// 									]
// 								}
// 							}`),
// 						},
// 					},
// 				},
// 			},
// 			want: want{
// 				rsp: &fnv1.RunFunctionResponse{
// 					Meta: &fnv1.ResponseMeta{Ttl: durationpb.New(60 * time.Second)},
// 					Desired: &fnv1.State{
// 						Resources: map[string]*fnv1.Resource{
// 							"argocd": {
// 								Resource: resource.MustStructJSON(`{
// 									"apiVersion": "helm.crossplane.io/v1beta1",
// 									"kind": "Release",
// 									"metadata": {
// 										"generateName": "example-xr-",
// 										"labels": {
// 											"crossplane.io/composite": "example-xr"
// 										}
// 									},
// 									"spec": {
// 										"forProvider": {
// 											"chart": {
// 												"name": "argo-cd",
// 												"repository": "https://argoproj.github.io/argo-helm",
// 												"version": "5.46.3"
// 											},
// 											"namespace": "argocd",
// 											"set": [
// 												{
// 													"name": "configs.credentialTemplates.github-private-repo-libraryManagementBackConfig.sshPrivateKey",
// 													"valueFrom": {
// 														"secretKeyRef": {
// 															"key": "ssh-private-key",
// 															"name": "ssh-private-key",
// 															"namespace": "argocd"
// 														}
// 													}
// 												},
// 												{
// 													"name": "configs.credentialTemplates.github-private-repo-libraryManagementAppOfApps.sshPrivateKey",
// 													"valueFrom": {
// 														"secretKeyRef": {
// 															"key": "ssh-private-key",
// 															"name": "ssh-private-key",
// 															"namespace": "argocd"
// 														}
// 													}
// 												}
// 											],
// 											"values": {
// 												"configs": {
// 													"credentialTemplates": {
// 														"github-private-repo-libraryManagementAppOfApps": {
// 															"url": "git@github.com:achrefbenmbarek1/libraryManagementAppOfApps.git"
// 														},
// 														"github-private-repo-libraryManagementBackConfig": {
// 															"url": "git@github.com:achrefbenmbarek1/libraryManagementBackConfig.git"
// 														}
// 													},
// 													"repositories": {
// 														"github-private-repo-libraryManagementAppOfApps": {
// 															"credentialName": "achrefbenmbarek1",
// 															"name": "libraryManagementAppOfApps",
// 															"type": "git",
// 															"url": "git@github.com:achrefbenmbarek1/libraryManagementAppOfApps.git"
// 														},
// 														"github-private-repo-libraryManagementBackConfig": {
// 															"credentialName": "achrefbenmbarek1",
// 															"name": "libraryManagementBackConfig",
// 															"type": "git",
// 															"url": "git@github.com:achrefbenmbarek1/libraryManagementBackConfig.git"
// 														}
// 													},
// 													"secret": {
// 														"argocdServerAdminPassword": "$2b$12$8z/N/HakTkOAiGnA4y7aGuVrlEvzierbqFuM9b0eL9EgX51ylQDWq",
// 														"createSecret": true
// 													}
// 												},
// 												"server": {
// 													"extraArgs": ["--insecure"],
// 													"ingress": {
// 														"enabled": true,
// 														"hostname": "argocd.local",
// 														"ingressClassName": "nginx",
// 														"path": "/",
// 														"pathType": "Prefix",
// 														"tls": false
// 													},
// 													"service": {
// 														"type": "ClusterIP"
// 													}
// 												}
// 											}
// 										},
// 										"providerConfigRef": {
// 											"name": "helm-provider"
// 										}
// 									}
// 								}`),
// 							},
// 						},
// 					},
// 					Conditions: []*fnv1.Condition{
// 						{
// 							Type:   "FunctionSuccess",
// 							Status: fnv1.Status_STATUS_CONDITION_TRUE,
// 							Reason: "Success",
// 							Target: fnv1.Target_TARGET_COMPOSITE_AND_CLAIM.Enum(),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	for name, tc := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			f := &Function{log: logging.NewNopLogger()}
// 			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)
//
// 			// Ignore metadata differences and generated fields
// 			opts := []cmp.Option{
// 				protocmp.Transform(),
// 				cmpopts.IgnoreUnexported(fnv1.RunFunctionResponse{}),
// 				cmpopts.IgnoreUnexported(fnv1.Resource{}),
// 				cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
// 					return k == "generateName" || k == "ownerReferences" || k == "annotations"
// 				}),
// 				cmpopts.SortMaps(func(a, b string) bool { return a < b }),
// 				// Fix: Use a concrete type for sorting slices
// 				cmpopts.SortSlices(func(a, b *fnv1.Resource) bool { return a.GetResource().String() < b.GetResource().String() }),
// 			}
//
// 			if diff := cmp.Diff(tc.want.rsp, rsp, opts...); diff != "" {
// 				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
// 			}
//
// 			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
// 				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
// 			}
// 		})
// 	}
// }
