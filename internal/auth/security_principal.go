package auth

import (
	"strings"

	"github.com/pglet/pglet/internal/utils"
)

type SecurityPrincipal struct {
	UID          string   `json:"uid"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	AuthProvider string   `json:"authProvider"`
	Token        string   `json:"token"`
	Groups       []string `json:"groups"`
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

		if strings.Index(permission, "/") != -1 && p.Groups != nil && len(p.Groups) > 0 {
			// check group
			for _, group := range p.Groups {
				if strings.ToLower(permission) == strings.ToLower(group) {
					identityMatched = true
					break
				}
			}
		} else if strings.ToLower(permission) == strings.ToLower(p.Username) ||
			// check username/password
			strings.ToLower(permission) == strings.ToLower(p.Email) {
			identityMatched = true
		}

		if authTypeMatched && identityMatched {
			return true
		}
	}
	return false
}
