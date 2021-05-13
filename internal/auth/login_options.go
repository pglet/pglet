package auth

import (
	"strings"

	"github.com/pglet/pglet/internal/utils"
)

type LoginOptions struct {
	GitHubEnabled    bool `json:"gitHubEnabled"`
	GitHubGroupScope bool `json:"gitHubGroupScope"`
	AzureEnabled     bool `json:"azureEnabled"`
	AzureGroupScope  bool `json:"azureGroupScope"`
	GoogleEnabled    bool `json:"googleEnabled"`
	GoogleGroupScope bool `json:"googleGroupScope"`
}

func GetLoginOptions(permissions string) *LoginOptions {

	if permissions == "" {
		return nil
	}

	opts := &LoginOptions{}

	// parse permissions
	permList := utils.SplitAndTrim(permissions, ",")

	for _, permission := range permList {
		// check permission's auth type
		authType := ""
		colonIdx := strings.Index(permission, ":")
		if colonIdx != -1 {
			authType = strings.ToLower(permission[:colonIdx])
			permission = permission[colonIdx+1:]
		}

		opts.GitHubEnabled = opts.GitHubEnabled || authType == "" || authType == GitHubAuth
		opts.AzureEnabled = opts.AzureEnabled || authType == "" || authType == AzureAuth
		opts.GoogleEnabled = opts.GoogleEnabled || authType == "" || authType == GoogleAuth

		// check if the requested permission is a group
		if strings.Index(permission, "/") != -1 {
			opts.GitHubGroupScope = opts.GitHubGroupScope || (opts.GitHubEnabled && authType == "" || authType == GitHubAuth)
			opts.AzureGroupScope = opts.AzureGroupScope || (opts.AzureEnabled && authType == "" || authType == AzureAuth)
			opts.GoogleGroupScope = opts.GoogleGroupScope || (opts.GoogleEnabled && authType == "" || authType == GoogleAuth)
		}
	}

	return opts
}
