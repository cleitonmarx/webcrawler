package testutil

import "github.com/cleitonmarx/webcrawler/config"

//FetcherStub implements a config.Repository interface
type ConfigRepoStub struct {
}

func (c *ConfigRepoStub) GetSystemConfiguration() (config.Configuration, error) {
	return config.Configuration{
		CurrentEnvironment: "Local",
		Environments: []config.EnvironmentConfig{
			config.EnvironmentConfig{
				Name: "Local",
				HTTPServer: config.HTTPServerConfig{
					Address: "127.0.0.1",
				},
			},
		},
	}, nil
}
