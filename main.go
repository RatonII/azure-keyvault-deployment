package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	aauth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/satori/go.uuid"
	"log"
	"runtime"
)

func main() {
	var sub string
	var resourceGr string
	var p KeyvaultProps
	ctx := context.Background()

	runtime.GOMAXPROCS(4)
	authorizer, err := aauth.NewAuthorizerFromCLI()
	if err != nil {
		log.Fatalf("unable to create azure authorizer: %v\n", err)
	}

	conf := flag.String("conf", "", "The path to the config file where the settings for keyvault are stored")
	subscription := flag.String("subscription", "", "The name of the subscription for cosmosdb or function to get the secrets from")
	resourceGroup := flag.String("resource-group", "", "The name of the resource group for cosmosdb or function to get the secrets from")
	flag.Parse()

	if *subscription != "" {
		sub = *subscription
	} else {
		log.Fatalln("Please provide a subscription for your azure account: --subscription")
	}
	if *conf == "" {
		log.Fatalln("Please provide a configuration file keyvault for your azure account: --conf")
	}
	if *resourceGroup != "" {
		resourceGr = *resourceGroup
	} else {
		log.Fatalln("Please provide a  resource group for your azure account: --resource-group")
	}
	vaultprops := p.getConf(conf)
	tenantid, err := uuid.FromString(vaultprops.TenantID)

	vaultIpRules := []keyvault.IPRule{}
	for _, IpRule := range vaultprops.NetworkAcls.IPRules {
		vaultIpRules = append(vaultIpRules, keyvault.IPRule{Value: &IpRule} )
	}

	vaultVirtualNetworkRules := []keyvault.VirtualNetworkRule{}
	for _, VirtualNetworkRule := range vaultprops.NetworkAcls.VirtualNetworkRules {
		vaultVirtualNetworkRules = append(vaultVirtualNetworkRules,
										keyvault.VirtualNetworkRule{ID: &VirtualNetworkRule} )
	}


	vaultProperties := &keyvault.VaultProperties{
		TenantID: &tenantid,
		Sku: &keyvault.Sku{
			Family: &vaultprops.SkuFamily,
			Name:   keyvault.SkuName(vaultprops.SkuName),
		},
		AccessPolicies:               &[]keyvault.AccessPolicyEntry{},
		EnabledForDeployment:         &vaultprops.EnabledForDeployment,
		EnabledForDiskEncryption:     &vaultprops.EnabledForDiskEncryption,
		EnabledForTemplateDeployment: &vaultprops.EnabledForTemplateDeployment,
		EnableSoftDelete:             &vaultprops.EnableSoftDelete,
		SoftDeleteRetentionInDays:    &vaultprops.SoftDeleteRetentionInDays,
		EnableRbacAuthorization:      &vaultprops.EnableRbacAuthorization,
		CreateMode:                   keyvault.CreateModeDefault,
		NetworkAcls:                  &keyvault.NetworkRuleSet{
			Bypass:              keyvault.NetworkRuleBypassOptions(vaultprops.NetworkAcls.Bypass),
			DefaultAction:       keyvault.NetworkRuleAction(vaultprops.NetworkAcls.DefaultAction),
			IPRules:             &vaultIpRules,
			VirtualNetworkRules: &vaultVirtualNetworkRules,
		},
	}

	vaultPatchProperties := &keyvault.VaultPatchProperties{
		TenantID: &tenantid,
		Sku: &keyvault.Sku{
			Family: &vaultprops.SkuFamily,
			Name:   keyvault.SkuName(vaultprops.SkuName),
		},
		EnabledForDeployment:         &vaultprops.EnabledForDeployment,
		EnabledForDiskEncryption:     &vaultprops.EnabledForDiskEncryption,
		EnabledForTemplateDeployment: &vaultprops.EnabledForTemplateDeployment,
		EnableSoftDelete:             &vaultprops.EnableSoftDelete,
		SoftDeleteRetentionInDays:    &vaultprops.SoftDeleteRetentionInDays,
		EnableRbacAuthorization:      &vaultprops.EnableRbacAuthorization,
		CreateMode:                   keyvault.CreateModeDefault,
		NetworkAcls:                  &keyvault.NetworkRuleSet{
			Bypass:              keyvault.NetworkRuleBypassOptions(vaultprops.NetworkAcls.Bypass),
			DefaultAction:       keyvault.NetworkRuleAction(vaultprops.NetworkAcls.DefaultAction),
			IPRules:             &vaultIpRules,
			VirtualNetworkRules: &vaultVirtualNetworkRules,
		},
	}
	authorizer, err = aauth.NewAuthorizerFromCLI()
	vaultclient := keyvault.NewVaultsClient(sub)
	vaultclient.Authorizer = authorizer
	if err != nil {
		log.Fatalf("Authentication failed with error: %s\n",err)
	}

	top := int32(10)
	getvault, err := vaultclient.List(ctx, &top)
	if err != nil {
		log.Fatalf("Listing top 10 keyvaults from subscription failed with error: %s\n",err)
	}

	vaultsnames := []string{}
	for _, vaults := range getvault.Values() {
		vaultsnames = append(vaultsnames, *vaults.Name)
	}

	if len(vaultIpRules) == 0 {
		vaultProperties.NetworkAcls.IPRules = nil
	}
	if len(vaultVirtualNetworkRules) == 0 {
		vaultProperties.NetworkAcls.VirtualNetworkRules = nil
	}
	if vaultprops.EnablePurgeProtection == false {
		vaultProperties.EnablePurgeProtection = nil
		vaultPatchProperties.EnablePurgeProtection = nil
	} else {
		vaultProperties.EnablePurgeProtection = &vaultprops.EnablePurgeProtection
		vaultPatchProperties.EnablePurgeProtection = &vaultprops.EnablePurgeProtection
	}
	if Find(vaultsnames, vaultprops.VaultName) == true {
		vault, err := vaultclient.Update(ctx, resourceGr, vaultprops.VaultName, keyvault.VaultPatchParameters{
			Tags:       vaultprops.Tags,
			Properties: vaultPatchProperties,
		})
		if err != nil {
			log.Fatalf("Keyvault update failed with error: %s\n",err)
		}
		fmt.Printf("Keyvault %s was updated\n",*vault.Name)
	} else if Find(vaultsnames, vaultprops.VaultName) == false {
		_, err := vaultclient.CreateOrUpdate(ctx, resourceGr, vaultprops.VaultName, keyvault.VaultCreateOrUpdateParameters{
			Location:   &vaultprops.Location,
			Tags:       vaultprops.Tags,
			Properties: vaultProperties,
		})
		if err != nil {
			log.Fatalf("Keyvault create failed with error: %s\n",err)
		}
		fmt.Printf("Keyvault %s was created\n", vaultprops.VaultName)
	}
}

//subscription id : 396e7375-178b-47ed-9e8b-bb198a66b766

