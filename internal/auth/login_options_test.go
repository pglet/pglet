package auth

import (
	"reflect"
	"testing"

	"github.com/pglet/pglet/internal/utils"
)

func TestLoginOptions(t *testing.T) {

	var loginOptionTests = []struct {
		permissions string        // input
		expected    *LoginOptions // expected result
	}{
		{"", nil},
		{"*", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: false, AzureEnabled: true, AzureGroupScope: false}},
		{"*/*", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: true, AzureEnabled: true, AzureGroupScope: true}},
		{"github:*", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: false, AzureEnabled: false, AzureGroupScope: false}},
		{"github:pglet/developers", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: true, AzureEnabled: false, AzureGroupScope: false}},
		{"azure:*", &LoginOptions{GitHubEnabled: false, GitHubGroupScope: false, AzureEnabled: true, AzureGroupScope: false}},
		{"azure:*/*", &LoginOptions{GitHubEnabled: false, GitHubGroupScope: false, AzureEnabled: true, AzureGroupScope: true}},
		{"*, azure:*/*", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: false, AzureEnabled: true, AzureGroupScope: true}},
		{"azure:*, github:*/*", &LoginOptions{GitHubEnabled: true, GitHubGroupScope: true, AzureEnabled: true, AzureGroupScope: false}},
	}

	for _, tt := range loginOptionTests {
		actual := GetLoginOptions(tt.permissions)
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("GetLoginOptions(%s): expected %v, actual %v", tt.permissions,
				utils.ToJSONIndent(tt.expected), utils.ToJSONIndent(actual))
		}
	}
}
