# Automatix

Automatix is a Github bot that updates Go Module dependencies and creates Pull
Requests.

It's intended to either be run:
- regularly via a scheduling mechanism like Cron or a CI Job or
- as standalone application on a server with `PeriodicIntervalMin` set to `>0`.

It clones the Github repositories listed in it's `config.toml` file.
For each repository it downloads updates for the libraries listed in the
`go.mod` file.
If libraries were updated it creates a Github Pull-Request for the changes.

If an open Pull-Request from Automatix already exists, it does not check for Go
module updates. This is to prevent that Automatix creates tons of Pull-Requests
for the same changes.

If the parameter `PeriodicIntervalMin` in the `config.toml` file is `0`,
Automatix only runs 1x time.
If the parameter is set to a value `>0`, automatix checks for updates in this
interval.


## Runtime Dependencies
- Git command-line tools
- Golang 1.11


## Quickstart
1. Run `go get -u github.com/simplesurance/automatix` to install Automatix
2. Create a OAUTH Github key for Automatix at https://github.com/settings/tokens
3. Run `./automatix` to create an example config
4. Move `config.example.toml` to `config.toml` and configure it according to your
   needs
5. Run `./automatix`


## Issues
- I'm getting a `404 Not Found` error when Automatix retrieves information from
  Github or tries to do `git push`
  - The permissions for your OAUTH key maybe not sufficient.
    Ensure the `rep - Full control of private repositories` checkbox is checked.
  - If you updated your OAUTH key and used Automatix before, you have to delete
    the previously checked out repositories in `RepoBaseDir`.
- The `git commit` step fails because my git user identity is not set
  - See https://help.github.com/articles/setting-your-commit-email-address-in-git/
    alternatively the environment variables `GIT_COMMITTER_EMAIL`,
    `GIT_AUTHOR_EMAIL` and `GIT_AUTHOR_NAME` can be set.
