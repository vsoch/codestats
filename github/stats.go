package github

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/vsoch/codestats/config"
	"github.com/vsoch/codestats/utils"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// These can exist anywhere in the repository
func hasContributing(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Contributing", Pass: existsRecursive(clone, "CONTRIBUTING")}
}
func hasCodeOfConduct(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Code-Of-Conduct", Pass: existsRecursive(clone, "CODE_OF_CONDUCT")}
}
func hasAuthors(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Authors", Pass: existsRecursive(clone, "AUTHORS")}
}
func hasSupport(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Support.md", Pass: existsRecursive(clone, "SUPPORT")}
}
func hasIssueTemplate(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Issue-Template", Pass: existsRecursive(clone, "ISSUE_TEMPLATE")}
}
func hasPullRequestTemplate(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Pull-Request-Template", Pass: existsRecursive(clone, "PULL_REQUEST_TEMPLATE")}
}
func hasSecurity(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Security", Pass: existsRecursive(clone, "SECURITY")}
}
func hasFunding(repo Repository, clone string) Stat {
	return Stat{Name: "Has-Funding", Pass: existsRecursive(clone, "FUNDING")}
}

// Determine if a file exists anywhere in a repository, and any casing
func existsRecursive(root string, pattern string) bool {
	var exists bool
	regex, _ := regexp.Compile(strings.ToLower(pattern))
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && regex.MatchString(strings.ToLower(info.Name())) {
			exists = true

			// This allows us to end early
			return io.EOF
		}
		return nil
	})
	return exists
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
		names = []string{"has-codeowners", "has-maintainers", "has-github-actions", "has-circle", "has-travis", "has-pull-approve", "has-glide", "has-code-of-conduct", "has-contributing", "has-authors", "has-pull-request-template", "has-issue-template", "has-support", "has-funding", "has-security"}

	}
	for _, name := range names {
		switch name {
		case "has-codeowners":
			tests = append(tests, hasCodeowners)
		case "has-maintainers":
			tests = append(tests, hasMaintainers)
		case "has-github-actions":
			tests = append(tests, hasGitHubActions)
		case "has-circle":
			tests = append(tests, hasCircle)
		case "has-travis":
			tests = append(tests, hasTravis)
		case "has-pull-approve":
			tests = append(tests, hasPullApprove)
		case "has-glide":
			tests = append(tests, hasGlide)
		case "has-code-of-conduct":
			tests = append(tests, hasCodeOfConduct)
		case "has-contributing":
			tests = append(tests, hasContributing)
		case "has-authors":
			tests = append(tests, hasAuthors)
		case "has-pull-request-template":
			tests = append(tests, hasPullRequestTemplate)
		case "has-issue-template":
			tests = append(tests, hasIssueTemplate)
		case "has-support":
			tests = append(tests, hasSupport)
		case "has-funding":
			tests = append(tests, hasFunding)
		case "has-security":
			tests = append(tests, hasSecurity)

		default:
			fmt.Println("Warning, unrecognized test", name)
		}
	}
	return tests
}

// Assemble a list of the functions to extract stats for!
func GetRepoStats(repoName string, yamlfile string, metrics []string) RepoResult {
	repo := GetRepo(repoName)
	directory, err := ioutil.TempDir("", "codestats")
	utils.CheckIfError(err)
	defer os.RemoveAll(directory)
	return getRepoStats(repo, directory, yamlfile, metrics)
}

func getRepoStats(repo Repository, directory string, yamlFile string, choices []string) RepoResult {

	// If we have a config, load it to get test preferences
	if yamlFile != "" && len(choices) == 0 {
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

	// Clone repository to inspect further
	cloneDir := filepath.Join(directory, filepath.Base(repo.Name))

	// Clone at depth 1 to be faster
	_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:               repo.CloneURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             1,
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
func GetOrgStats(orgName string, pattern string, skipPattern string, yamlFile string, choices []string) []RepoResult {

	repos := GetOrgRepos(orgName)
	results := []RepoResult{}

	// Create temporary directory to work in
	directory, err := ioutil.TempDir("", "codestats")
	utils.CheckIfError(err)
	defer os.RemoveAll(directory)

	regex, _ := regexp.Compile(pattern)
	regexSkip, _ := regexp.Compile(skipPattern)

	for _, repo := range repos {

		// SKip pattern?
		if skipPattern != "" && regexSkip.MatchString(repo.Name) {
			continue
		}

		// Don't parse repos that don't match
		if pattern != "" && !regex.MatchString(repo.Name) {
			continue
		}
		results = append(results, getRepoStats(repo, directory, yamlFile, choices))
	}
	return results
}
