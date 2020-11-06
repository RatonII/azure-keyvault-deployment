package main

import (
	"github.com/satori/go.uuid"
)

type KeyvaultProps struct {
	Location 					string					`yaml:"Location"`
	VaultName 					 string					`yaml:"VaultName"`
	TenantID                     string					`yaml:"TenantID"`
	SkuName                      string					`yaml:"SkuName"`
	SkuFamily                    string					`yaml:"SkuFamily,omitempty"`
	AccessPolicies               []VaultAccessPolicy	`yaml:"AccessPolicies,omitempty"`
	EnabledForDeployment         bool					`yaml:"EnabledForDeployment"`
	EnabledForDiskEncryption     bool					`yaml:"EnabledForDiskEncryption"`
	EnabledForTemplateDeployment bool					`yaml:"EnabledForTemplateDeployment"`
	EnableSoftDelete             bool					`yaml:"EnableSoftDelete"`
	SoftDeleteRetentionInDays    int32					`yaml:"SoftDeleteRetentionInDays"`
	EnableRbacAuthorization      bool					`yaml:"EnableRbacAuthorization"`
	EnablePurgeProtection        bool					`yaml:"EnablePurgeProtection,omitempty"`
	NetworkAcls                  VaultNetworkRuleSet	`yaml:"NetworkAcls,omitempty"`
	Tags 						 map[string]*string		`yaml:"Tags,omitempty"`
}
type VaultAccessPolicy struct {
	ObjectID                string
	ApplicationID           uuid.UUID
	KeysPermissions         []string
	SecretsPermissions      []string
	CertificatesPermissions []string
	StoragePermissions      []string
}

type VaultNetworkRuleSet struct {
	Bypass              string		`yaml:"Bypass,omitempty"`
	DefaultAction       string		`yaml:"DefaultAction,omitempty"`
	IPRules             []string	`yaml:"IPRules,omitempty"`
	VirtualNetworkRules []string	`yaml:"VirtualNetworkRules,omitempty"`
}
