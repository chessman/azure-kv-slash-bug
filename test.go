package main

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

var (
	vaultURL = fmt.Sprintf("https://%s.vault.azure.net", "<vaultName>")
	keyName  = "<keyName>"
	clientID = "<clientID>"
)

func AuthFromDeviceFlow(clientID, resource string) (autorest.Authorizer, error) {
	oauthClient := &autorest.Client{}
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, "common")
	if err != nil {
		return nil, err
	}

	deviceCode, err := adal.InitiateDeviceAuth(oauthClient, *oauthConfig, clientID, resource)
	if err != nil {
		return nil, err
	}

	fmt.Println(*deviceCode.Message)

	token, err := adal.WaitForUserCompletion(oauthClient, deviceCode)
	if err != nil {
		return nil, err
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(
		*oauthConfig,
		clientID,
		resource,
		*token,
	)
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(spt), nil
}

func KVTest(resource string) error {
	authorizer, err := AuthFromDeviceFlow(clientID, resource)
	if err != nil {
		return err
	}

	kv := keyvault.New()
	kv.Authorizer = authorizer
	str := "abc"

	_, err = kv.Encrypt(vaultURL, keyName, "", keyvault.KeyOperationsParameters{
		Algorithm: keyvault.RSAOAEP,
		Value:     &str,
	})

	return err
}

func main() {
	//no error
	if err := KVTest("https://vault.azure.net"); err != nil {
		fmt.Println(err)
	}

	//error
	if err := KVTest(azure.PublicCloud.KeyVaultEndpoint); err != nil {
		fmt.Println(err)
	}
}
