// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=template.fn.crossplane.io
// +versionName=v1beta1
package v1beta1input

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// TODO: Add your input type here! It doesn't need to be called 'Input', you can
// rename it to anything you like.

// Input can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane

// ArgoRepoGeneratorInput represents the input structure for your function.
type ArgoRepoGeneratorInput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ArgoRepoGeneratorSpec `json:"spec"`
}

// ArgoRepoGeneratorSpec defines the specification for the ArgoCD repository generator input.
type ArgoRepoGeneratorSpec struct {
	ForProvider ForProviderSpec `json:"forProvider"`
}

// ForProviderSpec defines the provider-specific fields.
type ForProviderSpec struct {
	Chart             ChartSpec  `json:"chart"`
	Namespace         string     `json:"namespace"`
	Values            ValuesSpec `json:"values"`
	Set               []SetItem  `json:"set"`
	ProviderConfigRef NameRef    `json:"providerConfigRef"`
}

// ChartSpec specifies the Helm chart information.
type ChartSpec struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Version    string `json:"version"`
}

// ValuesSpec contains customizable values for the ArgoCD setup.
type ValuesSpec struct {
	Server  ServerSpec `json:"server"`
	Configs ConfigSpec `json:"configs"`
}

// ServerSpec contains server-related configuration.
type ServerSpec struct {
	Service   ServerServiceSpec `json:"service"`
	Ingress   IngressSpec       `json:"ingress"`
	ExtraArgs []string          `json:"extraArgs"`
}

// ServerServiceSpec defines the service-related details.
type ServerServiceSpec struct {
	Type string `json:"type"`
}

// IngressSpec defines ingress configuration.
type IngressSpec struct {
	Enabled          bool   `json:"enabled"`
	IngressClassName string `json:"ingressClassName"`
	Hostname         string `json:"hostname"`
	Path             string `json:"path"`
	PathType         string `json:"pathType"`
	TLS              bool   `json:"tls"`
}

// ConfigSpec defines additional configuration options.
type ConfigSpec struct {
	Secret SecretSpec `json:"secret"`
	// CredentialTemplates map[string]CredentialTemplate `json:"credentialTemplates"`
	// Repositories        map[string]RepositorySpec     `json:"repositories"`
}

// SecretSpec contains secret-related information.
type SecretSpec struct {
	CreateSecret              bool   `json:"createSecret"`
	ArgoCDServerAdminPassword string `json:"argocdServerAdminPassword"`
}

// SetItem defines a single set item for value overrides.
type SetItem struct {
	Name      string    `json:"name"`
	ValueFrom ValueFrom `json:"valueFrom"`
}

// ValueFrom contains value source information.
type ValueFrom struct {
	SecretKeyRef SecretKeyRef `json:"secretKeyRef"`
}

// SecretKeyRef references a key in a secret.
type SecretKeyRef struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

// NameRef refers to a named configuration.
type NameRef struct {
	Name string `json:"name"`
}
