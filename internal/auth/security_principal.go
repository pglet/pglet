package auth

import (
	"errors"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pglet/pglet/internal/utils"
	"golang.org/x/oauth2"
)

type SecurityPrincipal struct {
	UID          string   `json:"uid"`
	AuthProvider string   `json:"authProvider"`
	Token        string   `json:"token"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Groups       []string `json:"groups"`
}

func NewPrincipal(authProvider string, groupsEnabled bool) *SecurityPrincipal {

	uid, _ := utils.GenerateRandomString(16)

	p := &SecurityPrincipal{
		UID:          uid,
		AuthProvider: authProvider,
	}

	if groupsEnabled {
		p.Groups = make([]string, 0)
	}

	return p
}

func (p *SecurityPrincipal) SetToken(token *oauth2.Token) error {
	if token == nil {
		p.Token = ""
		return nil
	}

	j := utils.ToJSON(token)
	enc, err := utils.EncryptWithMasterKey([]byte(j))
	if err != nil {
		return err
	}
	p.Token = string(enc)
	return nil
}

func (p *SecurityPrincipal) GetToken() (*oauth2.Token, error) {
	if p.Token == "" {
		return nil, nil
	}

	j, err := utils.DecryptWithMasterKey([]byte(p.Token))
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	utils.FromJSON(string(j), token)
	return token, nil
}

func (p *SecurityPrincipal) UpdateDetails() error {
	if p.AuthProvider == GitHubAuth {
		return p.updateFromGitHub()
	} else if p.AuthProvider == AzureAuth {
		return p.updateFromAzure()
	} else if p.AuthProvider == "" {
		return errors.New("auth provider is not set")
	} else {
		return errors.New("unknown auth provider")
	}
}

func (p *SecurityPrincipal) updateFromGitHub() error {
	return nil
}

func (p *SecurityPrincipal) updateFromAzure() error {
	return nil
}

func (p *SecurityPrincipal) HasPermissions(permissions string) bool {

	if permissions == "" {
		return true
	}

	if p.AuthProvider == "" {
		return false
	}

	// parse permissions
	permList := utils.SplitAndTrim(permissions, ",")

	for _, permission := range permList {

		// check permission's auth type
		authType := ""
		colonIdx := strings.Index(permission, ":")
		if colonIdx != -1 {
			authType = strings.ToLower(strings.TrimSpace(permission[:colonIdx]))
			permission = strings.TrimSpace(permission[colonIdx+1:])
		}

		authTypeMatched := authType == "" || p.AuthProvider == authType
		identityMatched := false

		pg := glob.MustCompile(strings.ToLower(permission))

		if strings.Index(permission, "/") != -1 && p.Groups != nil && len(p.Groups) > 0 {
			// check group
			for _, group := range p.Groups {
				if pg.Match(strings.ToLower(group)) {
					identityMatched = true
					break
				}
			}
		} else if (p.Username != "" && pg.Match(strings.ToLower(p.Username))) ||
			(p.Email != "" && pg.Match(strings.ToLower(p.Email))) {
			identityMatched = true
		}

		if authTypeMatched && identityMatched {
			return true
		}
	}
	return false
}
