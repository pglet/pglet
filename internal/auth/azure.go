package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pglet/pglet/internal/auth/msgraph"
	"github.com/pglet/pglet/internal/utils"
	"golang.org/x/oauth2"
)

func (p *SecurityPrincipal) updateFromAzure() error {

	// retrieve token
	token, err := p.GetToken()
	if err != nil {
		return err
	}

	if token == nil {
		return errors.New("Azure OAuth token is not set")
	}

	// GitHub client
	oauthConfig := GetOauthConfig(p.AuthProvider, p.Groups != nil)
	client := oauthConfig.Client(oauth2.NoContext, token)

	// get user details
	// https://docs.microsoft.com/en-us/graph/api/user-get?view=graph-rest-1.0&tabs=http
	graphUser := &msgraph.User{}
	err = GetObject(client, "https://graph.microsoft.com/v1.0/me", graphUser)
	if err != nil {
		return err
	}

	p.Login = graphUser.UserPrincipalName
	p.Name = graphUser.DisplayName
	if graphUser.Mail == "" {
		p.Email = graphUser.UserPrincipalName
	}

	if p.Groups != nil {
		p.Groups = make([]string, 0)

		// read organization details
		// https://docs.microsoft.com/en-us/graph/api/organization-get?view=graph-rest-1.0&tabs=http
		orgResp := &msgraph.OrganizationsResponse{}
		err = GetObject(client, "https://graph.microsoft.com/v1.0/organization", orgResp)
		if err != nil {
			return err
		}

		if orgResp.Values == nil || len(orgResp.Values) == 0 {
			return nil
		}

		org := orgResp.Values[0] // there is an array with only one org

		// read user memberships
		// https://docs.microsoft.com/en-us/graph/api/user-list-memberof?view=graph-rest-1.0&tabs=http
		groupsResp := &msgraph.MembershipsResponse{}
		err = GetObject(client, "https://graph.microsoft.com/v1.0/me/memberOf?$select=displayName,id", groupsResp)
		if err != nil {
			return err
		}

		for _, group := range groupsResp.Values {
			p.Groups = append(p.Groups, fmt.Sprintf("%s/%s", org.Id, group.DisplayName))
		}
	}

	return nil
}

func GetObject(client *http.Client, url string, dst interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	utils.FromJSON(string(body), dst)

	return nil
}
