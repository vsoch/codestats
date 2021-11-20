package github

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/vsoch/org-stats/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// A stat can either pass or fail
type Stat struct {
	Name string
	Pass bool
}

// RepoResult has a list of subfields we care about
type RepoResult struct {
	Stats     []Stat
	Name      string
	Branch    string
	Url       string
	Stars     int
	Forks     int
	Issues    int
	Language  string
	Archived  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Every defined stat needs to return a Stat object
func hasCodeowners(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Codeowners", Pass: exists(clone, "CODEOWNERS")}
}
func hasMaintainers(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Maintainers", Pass: exists(clone, "MAINTAINERS")}
}
func hasGitHubActions(repo Repository, clone string) Stat {
	return Stat{Name: "Has-GitHub-Actions", Pass: exists(clone, ".github/workflows")}
}
func hasCircle(repo Repository, clone string) Stat {
	return Stat{Name: "Has-CircleCI", Pass: exists(clone, ".circleci")}
}
func hasTravis(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Travis", Pass: exists(clone, ".travis.yml")}
}
func hasPullApprove(repo Repository, clone string) Stat {
	return Stat{Name: "Has-PullApprove", Pass: exists(clone, ".pullapprove.yml")}
}
func hasGlide(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Glide", Pass: exists(clone, ".glide.yaml")}
}

// hasFile is a general function to test if a repository root has a file
func exists(dirname string, filename string) bool {
	path := filepath.Join(dirname, filename)
	return utils.Exists(path)
}

// getTests for different stats depending on the user config
func getTests() []func(repo Repository, clone string) Stat {
	var tests []func(repo Repository, clone string) Stat
	tests = append(tests, hasCodeowners)
	tests = append(tests, hasMaintainers)
	tests = append(tests, hasGitHubActions)
	tests = append(tests, hasCircle)
	tests = append(tests, hasTravis)
	tests = append(tests, hasPullApprove)
	tests = append(tests, hasGlide)
	return tests
}

// Assemble a list of the functions to extract stats for!
func GetOrgStats(orgName string) []RepoResult {

	repos := GetOrgRepos(orgName)
	results := []RepoResult{}

	// Create temporary directory to work in
	directory, err := ioutil.TempDir("", "org-stats")
	utils.CheckIfError(err)
	defer os.RemoveAll(directory)

	for _, repo := range repos {

		language := repo.Language
		if language == "null" {
			language = "?"
		}
		result := RepoResult{Name: repo.Name, Stars: repo.StargazersCount,
			Forks: repo.ForksCount, Issues: repo.OpenIssuesCount,
			Language: repo.Language, Branch: repo.DefaultBranch,
			Archived: repo.Archived, CreatedAt: repo.CreatedAt,
			UpdatedAt: repo.UpdatedAt}

		fmt.Println(repo.Name)

		// Clone repository to inspect further
		cloneDir := filepath.Join(directory, filepath.Base(repo.Name))
		r, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
			URL:               repo.CloneURL,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		utils.CheckIfError(err)
		fmt.Println(r)

		// Test results
		stats := []Stat{}

		// Load tests
		for _, test := range getTests() {
			stats = append(stats, test(repo, cloneDir))
		}

		result.Stats = stats
		results = append(results, result)
	}
	return results
}
