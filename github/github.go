package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vsoch/codestats/utils"
)

func githubGetRequest(url string) string {
	headers := make(map[string]string)
	token := os.Getenv("GITHUB_TOKEN")
	headers["Accept"] = "application/vnd.github.v3+json"
	if token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", token)
	}
	return utils.GetRequest(url, headers)
}

func GetReleases(name string) Releases {

	response := githubGetRequest("https://api.github.com/repos/" + name + "/releases")

	// The response gets parsed into a spack package
	releases := Releases{}
	err := json.Unmarshal([]byte(response), &releases)
	if err != nil {
		log.Fatalf("Issue unmarshalling releases data structure\n")
	}
	return releases
}

func GetOrgRepos(orgName string) Repos {

	response := githubGetRequest("https://api.github.com/orgs/" + orgName + "/repos")

	// The response gets parsed into a spack package
	repos := Repos{}
	err := json.Unmarshal([]byte(response), &repos)
	if err != nil {
		log.Fatalf("Issue unmarshalling repositories data structure\n")
	}
	return repos
}

func GetRepo(repoName string) Repository {

	response := githubGetRequest("https://api.github.com/repos/" + repoName)

	// The response gets parsed into a spack package
	repo := Repository{}
	err := json.Unmarshal([]byte(response), &repo)
	if err != nil {
		log.Fatalf("Issue unmarshalling repository data structure\n")
	}
	return repo
}

func GetCommits(name string, branch string) Commits {
	url := "https://api.github.com/repos/" + name + "/commits"

	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.github.v3+json"
	headers["Sha"] = branch
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		headers["Authorization"] = fmt.Sprintf("token %s", token)
	}
	response := utils.GetRequest(url, headers)

	commits := Commits{}
	err := json.Unmarshal([]byte(response), &commits)
	if err != nil {
		log.Fatalf("Issue unmarshalling commits data structure\n")
	}
	return commits
}
