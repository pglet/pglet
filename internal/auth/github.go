package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (p *SecurityPrincipal) updateFromGitHub() error {

	// retrieve token
	token, err := p.GetToken()
	if err != nil {
		return err
	}

	if token == nil {
		return errors.New("GitHub OAuth token is not set")
	}

	// GitHub client
	oauthConfig := GetOauthConfig(p.AuthProvider, p.Groups != nil)
	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, token))

	// read user details
	githubUser, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return err
	}

	p.ID = strconv.FormatInt(githubUser.GetID(), 10)
	p.Login = githubUser.GetLogin()
	p.Name = githubUser.GetName()

	// read user emails
	listEmailOpts := &github.ListOptions{
		PerPage: 10,
	}
	for {
		emails, resp, err := client.Users.ListEmails(context.Background(), listEmailOpts)
		if err != nil {
			return err
		}

		for _, email := range emails {
			if *email.Primary {
				p.Email = *email.Email
				break
			}
		}

		if p.Email != "" {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		listEmailOpts.Page = resp.NextPage
	}

	// read user teams
	if p.Groups != nil {
		p.Groups = make([]string, 0)

		listTeamsOpts := &github.ListOptions{
			PerPage: 10,
		}
		for {
			teams, resp, err := client.Teams.ListUserTeams(context.Background(), listTeamsOpts)
			if err != nil {
				return err
			}

			for _, team := range teams {
				p.Groups = append(p.Groups, fmt.Sprintf("%s/%s", *team.Organization.Login, team.GetName()))
			}

			if resp.NextPage == 0 {
				break
			}
			listTeamsOpts.Page = resp.NextPage
		}
	}
	return nil
}
