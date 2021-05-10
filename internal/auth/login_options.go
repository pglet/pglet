package auth

type LoginOptions struct {
	GitHubEnabled    bool `json:"gitHubEnabled"`
	GitHubGroupScope bool `json:"gitHubGroupScope"`
	AzureEnabled     bool `json:"azureEnabled"`
	AzureGroupScope  bool `json:"azureGroupScope"`
}
