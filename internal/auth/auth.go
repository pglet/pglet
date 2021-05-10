package auth

const (
	GitHubAuth = "github"
	AzureAuth  = "azure"

	// Users with both a personal Microsoft account and a work or school account from Azure AD can sign in to the application.
	AzureCommonTenant = "common"

	// Only users with work or school accounts from Azure AD can sign in to the application.
	AzureOrganizationsTenant = "organizations"

	// Only users with a personal Microsoft account can sign in to the application.
	AzureConsumersTenant = "consumers"

	/*
		8eaef023-2b34-4da1-9baa-8bc8c9d6a490 or contoso.onmicrosoft.com

		Only users from a specific Azure AD tenant (whether they are members in the directory with a work or school account,
		or they are guests in the directory with a personal Microsoft account) can sign in to the application.
		Either the friendly domain name of the Azure AD tenant or the tenant's GUID identifier can be used.
		You can also use the consumer tenant, 9188040d-6c67-4c5b-b112-36a304b66dad, in place of the consumers tenant.

		More info: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-protocols-oidc
	*/
)
