// +build integration

package az

import (
	"context"
	"math/rand"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// setup parameters needed to run integration tests
func setup(t *testing.T) (name, subscriptionID, tenantID string) {
	t.Log("setup integration test variables")
	tenantID = os.Getenv("TENANT_ID")
	if len(tenantID) == 0 {
		t.Fatalf("expected non-empty TENANT_ID environment variable")
	}

	subscriptionID = os.Getenv("SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatalf("expected non-empty SUBSCRIPTION_ID environment variable")
	}

	// TODO: Check possibility of using https://github.com/google/uuid instead of generateRandomString()
	name = "epiphany-cli-tests-" + generateRandomString(6)
	t.Logf("integration tests will use name %s as created application name", name)
	return
}

func TestCreateServicePrincipal(t *testing.T) {
	//given
	name, subscriptionID, tenantID := setup(t)
	pass := generateRandomString(10)
	env, err := azure.EnvironmentFromName(cloudName)
	if err != nil {
		t.Fatal(err)
	}
	authorizer, err := auth.NewAuthorizerFromCLIWithResource(env.GraphEndpoint)
	if err != nil {
		t.Fatal(err)
	}
	spClient := getTestServicePrincipalClient(tenantID, authorizer)
	appClient := getTestApplicationClient(tenantID, authorizer)

	// when
	var spObjectID, appObjectID string
	app, sp, err := CreateServicePrincipal(pass, subscriptionID, tenantID, name)
	defer cleanupTestResources(spObjectID, appObjectID, spClient, appClient, t)
	if sp != nil {
		spObjectID = *sp.ObjectID
	}
	if app != nil {
		appObjectID = *app.ObjectID
	}
	if err != nil {
		if appObjectID == "" && spObjectID == "" {
			t.Fatal(err)
		} else {
			// appID and spID were returned so there is potential to cleanup
			t.Error(err)
		}
	}

	// then
	servicePrincipal, err := spClient.Get(context.TODO(), spObjectID)
	if err != nil {
		t.Error(err)
	}
	if servicePrincipal.Response.Status != "200 OK" {
		t.Errorf("service principal GET operation returned: %s", servicePrincipal.Response.Status)
	}

	application, err := appClient.Get(context.TODO(), appObjectID)
	if err != nil {
		t.Error(err)
	}
	if application.Response.Status != "200 OK" {
		t.Errorf("application GET operation returned: %s", servicePrincipal.Response.Status)
	}

	credentials := Credentials{
		AppID:          appObjectID,
		Password:       pass,
		Tenant:         tenantID,
		SubscriptionID: subscriptionID,
	}

	t.Logf("created credentials: %#v", credentials)
}

// cleanupTestResources cleans up Service Principal and related resources based on app and sp object ID
func cleanupTestResources(spObjectID, appObjectID string, servicePrincipalClient graphrbac.ServicePrincipalsClient, applicationsClient graphrbac.ApplicationsClient, t *testing.T) {
	t.Log("start deleting service principal.")

	_, err := servicePrincipalClient.Delete(context.TODO(), spObjectID)
	if err != nil {
		t.Error(err)
	}

	_, err = applicationsClient.Delete(context.TODO(), appObjectID)
	if err != nil {
		t.Error(err)
	}
}

// getTestServicePrincipalClient gets service principal client for test purposes
func getTestServicePrincipalClient(tenantID string, authorizer autorest.Authorizer) graphrbac.ServicePrincipalsClient {
	spClient := graphrbac.NewServicePrincipalsClient(tenantID)
	spClient.Authorizer = authorizer

	return spClient
}

// getTestApplicationClient gets application client for test purposes
func getTestApplicationClient(tenantID string, authorizer autorest.Authorizer) graphrbac.ApplicationsClient {
	appClient := graphrbac.NewApplicationsClient(tenantID)
	appClient.Authorizer = authorizer

	return appClient
}

func generateRandomString(length int) string {
	hexBytes := "1234567890abcdef"
	b := make([]byte, length)
	for i := range b {
		b[i] = hexBytes[rand.Intn(len(hexBytes))]
	}
	return string(b)
}
