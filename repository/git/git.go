package git

import (
	"fmt"
	"path"
	"strings"

	"github.com/simplesurance/automatix/exec"
	"github.com/simplesurance/automatix/fs"

	"github.com/pkg/errors"
)

// Repository represents a git repository
type Repository struct {
	dir string
	url string
}

// NewGithubRepository initializes a github based repository object
func NewGithubRepository(oauthToken, organization, repository, localdirectory string) *Repository {
	return &Repository{
		dir: localdirectory,
		url: GithubRepoURL(oauthToken, organization, repository),
	}
}

// Clone clones the git repository, the directory must not exist
func (r *Repository) Clone(branch string) error {
	cmd := fmt.Sprintf("git clone -b %s %s %s", branch, r.url, r.dir)
	_, err := exec.Command("", cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	return nil
}

// Update checksout the branch and updates it to the newest remote ref
func (r *Repository) Update(branch string) error {
	cmd := fmt.Sprintf("git fetch origin %s", branch)
	_, err := exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	cmd = fmt.Sprintf("git checkout -f %s", branch)
	_, err = exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	cmd = fmt.Sprintf("git reset --hard origin/%s", branch)
	_, err = exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	return nil
}

// UpClone clones the repository if the directory does not exist yet otherwise
// updates it.
func (r *Repository) UpClone(branch string) error {
	if fs.PathExists(r.dir) {
		if !fs.IsDir(path.Join(r.dir, ".git")) {
			return fmt.Errorf("%q exist but is not a git directory", r.dir)
		}

		err := r.Update(branch)
		if err != nil {
			return errors.Wrap(err, "updating repository failed")
		}

		return nil
	}

	err := r.Clone(branch)
	if err != nil {
		return errors.Wrap(err, "cloning repository failed")
	}

	return nil
}

// GithubRepoURL returns a HTTPS url to acces the github repository with OAUTH
// authentication
func GithubRepoURL(oauthToken, organization, repository string) string {
	return fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s", oauthToken, organization, repository)
}

// FileChanged returns true if the tracked filed contains non-committed changes
func (r *Repository) FileChanged(relpath string) (bool, error) {
	fname := path.Base(relpath)
	cmd := fmt.Sprintf("git ls-files -m %s", relpath)

	out, err := exec.Command(r.dir, cmd)
	if err != nil {
		return false, errors.Wrap(err, "running git clone failed")
	}

	if len(out) < len(fname) {
		return false, nil
	}

	return true, nil
}

// CreateBranch creates a new branch
func (r *Repository) CreateBranch(branch string) error {
	cmd := fmt.Sprintf("git checkout -t -b %s", branch)
	_, err := exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	return nil
}

// CreateCommit creates a commit
func (r *Repository) CreateCommit(msg string, files []string) error {
	cmd := fmt.Sprintf("git commit -m \"%s\" %s", msg, strings.Join(files, " "))
	_, err := exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)

	}

	return nil
}

// Push pushes the branch to the default remote
func (r *Repository) Push() error {
	cmd := "git push origin"
	_, err := exec.Command(r.dir, cmd)
	if err != nil {
		return errors.Wrapf(err, "running %q failed", cmd)
	}

	return nil
}
