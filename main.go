package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"

	"github.com/simplesurance/automatix/fs"
	"github.com/simplesurance/automatix/log"
	"github.com/simplesurance/automatix/repository/git"
	"github.com/simplesurance/automatix/vendortool/gomod"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	homepage = "https://github.com/simplesurance/automatix"
	toolName = "Automatix"

	cfgPath        = "config.toml"
	cfgExamplePath = "config.example.toml"

	githubPRTitle = "Go Module Update by " + toolName
	commitMsg     = "go: update vendor libraries\nUpdate was done by " + toolName + "(" + homepage + ")"
)

func prExists(clt *github.Client, owner, repository, title string) (bool, string, error) {
	page := 0

	for {
		prs, resp, err := clt.PullRequests.List(
			context.Background(), owner, repository,
			&github.PullRequestListOptions{
				State: "open",
				ListOptions: github.ListOptions{
					PerPage: 10,
					Page:    page,
				}})
		if err != nil {
			return false, "", errors.Wrap(err, "fetching pull requests from github failed")
		}

		for _, pr := range prs {
			if *pr.Title == githubPRTitle {
				return true, *pr.HTMLURL, nil
			}

		}

		if resp.NextPage == 0 {
			return false, "", nil
		}

		page = resp.NextPage
	}
}

func newGithubClient(cfg *Config) *github.Client {
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.GitHubToken})
	tokenClt := oauth2.NewClient(context.Background(), tokenSrc)

	return github.NewClient(tokenClt)
}

func updateVendorDeps(githubClt *github.Client, cfg *Config, repo *GitHubRepository) bool {
	var depsUpdated bool

	ts := time.Now()
	branchName := fmt.Sprintf("%s-%d", toolName, ts.Unix())
	ldir := path.Join(cfg.RepoBaseDir, repo.Repository)
	name := path.Base(repo.Repository)

	exists, prURL, err := prExists(githubClt, repo.Owner, repo.Repository, githubPRTitle)
	if err != nil {
		log.Errorf("%s: finding pr failed: %s", name, err)
		return false
	}

	if exists {
		log.Infof("%s: open PR already exists (%s), skipping repository check\n", name, prURL)
		return true
	}

	gitRepo := git.NewGithubRepository(cfg.GitHubToken, repo.Owner, repo.Repository, ldir)

	log.Infof("%s: fetching %s branch", name, repo.Branch)
	err = gitRepo.UpClone(repo.Branch)
	if err != nil {
		log.Errorf("%s: fetching newest ref of branch %q failed: %s", name, repo.Branch, err)
		return false
	}

	goModPath := path.Join(ldir, repo.GoModPath)
	if !fs.IsFile(goModPath) {
		log.Errorf("%s: go.Mod file %s does not exist in repository", name, repo.GoModPath)
		return false
	}

	log.Infof("%s: updating vendor dependencies", name)
	mod := gomod.NewModule(path.Join(ldir, repo.GoModPath))
	if repo.MinorVersionUpdate {
		depsUpdated, err = mod.FullUpdate()
	} else {
		depsUpdated, err = mod.SafeUpdate()
	}
	if err != nil {
		log.Errorf("%s: updating vendor dependencies failed: %s", name, err.Error())
		return false
	}

	if !depsUpdated {
		log.Infof("%s: vendor dependencies already uptodate", name)
		return false
	}

	log.Infof("%s: vendor dependencies updated", name)

	if repo.VendorLibs {
		log.Infof("%s downloading dependencies into vendor directory", name)
		err := mod.VendorDependencies()
		if err != nil {
			log.Errorf("%s: download vendor dependencies failed: %s", name, err)
			return false
		}
	}

	log.Infof("%s: creating and pushing new git branch %q", name, branchName)
	err = gitRepo.CreateBranch(branchName)
	if err != nil {
		log.Errorf("%s: creating git branch failed: %s", name, err)
		return false
	}

	err = gitRepo.CreateCommit(commitMsg, mod.Files())
	if err != nil {
		log.Errorf("creating git branch failed: %s", err.Error())
		return false
	}

	log.Infof("%s: pushing git branch", name)
	err = gitRepo.Push()
	if err != nil {
		log.Errorf("%s: pushing git branch failed: %s", name, err)
		return false
	}

	prTitle := githubPRTitle
	prBody := commitMsg
	log.Infof("%s: creating github PR %q", name, prTitle)
	pr, _, err := githubClt.PullRequests.Create(context.Background(), repo.Owner,
		repo.Repository, &github.NewPullRequest{
			Base:  &repo.Branch,
			Head:  &branchName,
			Title: &prTitle,
			Body:  &prBody,
		})

	if err != nil {
		log.Errorf("%s: creating github PR failed: %s", name, err)
		return false
	}

	log.Infof("%s: pull request created (%s)", name, *pr.HTMLURL)

	return true
}

func main() {
	success := true

	if !fs.IsFile(cfgPath) {
		log.Infof("config file %q does not exist", cfgPath)

		err := ExampleConfig.ToFile(cfgExamplePath)
		if err != nil {
			log.Fatalf("writing example config file to %q failed: %s", cfgExamplePath, err)
		}
		log.Infof("written example config file to %q", cfgExamplePath)

		os.Exit(1)
	}

	cfg, err := ConfigFromFile(cfgPath)
	if err != nil {
		log.Fatalf("reading configuration file %q failed: %s", cfgPath, err)
	}

	log.Infof("configuration loaded from %q, configuration:\n%s\n", cfgPath, cfg)

	fs.Mkdir(cfg.RepoBaseDir)

	if len(cfg.GitHubToken) == 0 {
		cfg.GitHubToken = os.Getenv("GITHUB_OAUTH_TOKEN")
	}

	githubClt := newGithubClient(cfg)

	for {
		log.Infof("--")
		for _, repo := range cfg.Repositories {
			if !updateVendorDeps(githubClt, cfg, repo) {
				success = false
			}
		}

		if cfg.PeriodicIntervalMin == 0 {
			break
		}

		log.Infof("--")
		log.Infof("Next module update check in %vmin", cfg.PeriodicIntervalMin)
		time.Sleep(time.Duration(cfg.PeriodicIntervalMin) * time.Minute)
	}

	if !success {
		os.Exit(1)
	}

	os.Exit(0)
}
