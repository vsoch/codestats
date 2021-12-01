package github

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/vsoch/codestats/config"
	"github.com/vsoch/codestats/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
func getTests(names []string) []func(repo Repository, clone string) Stat {
	var tests []func(repo Repository, clone string) Stat

	// If names is empty, assume all
	if len(names) == 0 {
		names = []string{"hasCodeowners", "hasMaintainers", "hasGitHubActions", "hasCircle", "hasTravis", "hasPullApprove", "hasGlide"}
	}
	for _, name := range names {
		switch name {
		case "hasCodeowners":
			tests = append(tests, hasCodeowners)
		case "hasMaintainers":
			tests = append(tests, hasMaintainers)
		case "hasGitHubActions":
			tests = append(tests, hasGitHubActions)
		case "hasCircle":
			tests = append(tests, hasCircle)
		case "hasTravis":
			tests = append(tests, hasTravis)
		case "hasPullApprove":
			tests = append(tests, hasPullApprove)
		case "hasGlide":
			tests = append(tests, hasGlide)
		default:
			fmt.Println("Warning, unrecognized test %s", name)
		}
	}
	return tests
}

// Assemble a list of the functions to extract stats for!
func GetRepoStats(repoName string, yamlfile string) RepoResult {
	repo := GetRepo(repoName)
	directory, err := ioutil.TempDir("", "codestats")
	utils.CheckIfError(err)
	defer os.RemoveAll(directory)
	return getRepoStats(repo, directory, yamlfile)
}

func getRepoStats(repo Repository, directory string, yamlFile string) RepoResult {

	// If we have a config, load it to get test preferences
	choices := []string{}
	if yamlFile != "" {
		conf := config.Load(yamlFile)
		choices = conf.Stats
	}

	language := repo.Language
	if language == "null" {
		language = "?"
	}
	result := RepoResult{Name: repo.FullName, Stars: repo.StargazersCount,
		Forks: repo.ForksCount, Issues: repo.OpenIssuesCount,
		Language: repo.Language, Branch: repo.DefaultBranch,
		Archived: repo.Archived, CreatedAt: repo.CreatedAt,
		UpdatedAt: repo.UpdatedAt, Url: repo.HTMLURL}

	fmt.Println(repo.Name)
	fmt.Println()

	// Clone repository to inspect further
	cloneDir := filepath.Join(directory, filepath.Base(repo.Name))

	_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:               repo.CloneURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	utils.CheckIfError(err)

	// Test results
	stats := []Stat{}

	// Load tests based on user choices
	for _, test := range getTests(choices) {
		stats = append(stats, test(repo, cloneDir))
	}

	result.Stats = stats
	return result
}

// Assemble a list of the functions to extract stats for!
func GetOrgStats(orgName string, pattern string, yamlFile string) []RepoResult {

	repos := GetOrgRepos(orgName)
	results := []RepoResult{}

	// Create temporary directory to work in
	directory, err := ioutil.TempDir("", "codestats")
	utils.CheckIfError(err)
	defer os.RemoveAll(directory)

	regex, _ := regexp.Compile(pattern)

	for _, repo := range repos {

		// Don't parse repos that don't match
		if pattern != "" && !regex.MatchString(repo.Name) {
			continue
		}
		results = append(results, getRepoStats(repo, directory, yamlFile))
	}
	return results
}
