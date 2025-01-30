// Package models contains data structures for repository configurations
package models

// Repository represents a Git repository configuration.
type Repository struct {
	URL            string `json:"url" yaml:"url"`
	Type           string `json:"type" yaml:"type"`
	Name           string `json:"name" yaml:"name"`
	CredentialName string `json:"credentialName" yaml:"credentialName"`
}
