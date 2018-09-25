package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// Config represents the Automatix configuration
type Config struct {
	RepoBaseDir         string
	Repositories        []*GitHubRepository
	GitHubToken         string `comment:"github oauth token, if empty the token from the environment variable GITHUB_OAUTH_TOKEN is used"`
	PeriodicIntervalMin uint32 `comment:"set to 0 to run update only 1x"`
}

// GitHubRepository represents a repository configuration
type GitHubRepository struct {
	Owner              string
	Repository         string
	Branch             string
	VendorLibs         bool   `comment:"download dependencies into vendor/ directory after update"`
	MinorVersionUpdate bool   `comment:"if enabled libraries will be updated to newest minor or patch version, if false only to newest patch version."`
	GoModPath          string `comment:"path to go.mod file in repository"`
}

// ExampleConfig represents the example automatix config file
var ExampleConfig = Config{
	RepoBaseDir:         fmt.Sprintf("/var/lib/%s/repositories/", strings.ToLower(toolName)),
	GitHubToken:         "1234567890",
	PeriodicIntervalMin: 0,
	Repositories: []*GitHubRepository{
		&GitHubRepository{
			Owner:              "simplesurance",
			Repository:         "baur",
			Branch:             "master",
			VendorLibs:         true,
			MinorVersionUpdate: true,
			GoModPath:          "go.mod",
		},

		&GitHubRepository{
			Owner:              "simplesurance",
			Repository:         "automatix",
			Branch:             "master",
			VendorLibs:         true,
			MinorVersionUpdate: false,
			GoModPath:          "go.mod",
		},
	},
}

func (c *Config) String() string {
	res := fmt.Sprintf("RepoBaseDir: %v\nPeriodicIntervalMin: %v\n",
		c.RepoBaseDir, c.PeriodicIntervalMin)

	for i, repo := range c.Repositories {
		res += repo.String()
		if i+1 < len(c.Repositories) {
			res += "--\n"
		}
	}

	return res
}

// ToFile writes the Config object to a file
func (c *Config) ToFile(path string) error {
	data, err := toml.Marshal(*c)
	if err != nil {
		return errors.Wrap(err, "marshalling config failed")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "opening file failed")
	}

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrap(err, "writing to file failed")
	}

	err = f.Close()
	if err != nil {
		return errors.Wrap(err, "closing file failed")
	}

	return err
}

// ConfigFromFile reads a configuration file and returns it as Config object
func ConfigFromFile(path string) (*Config, error) {
	var config Config

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file failed")
	}

	err = toml.Unmarshal(content, &config)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling file")
	}

	return &config, nil
}

func (g *GitHubRepository) String() string {
	return fmt.Sprintf(" Owner: %v\n Repository: %v\n Branch: %v\n VendorLibs: %v\n MinorVersionUpdate: %v",
		g.Owner, g.Repository, g.Branch, g.VendorLibs, g.MinorVersionUpdate)
}
