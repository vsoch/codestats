package cli

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/vsoch/codestats/github"
)

// Args and flags
type RepoArgs struct {
	Repos []string `desc:"One or more GitHub repository names to parse."`
}
type RepoFlags struct {
	Pretty bool `long:"pretty" desc:"If printing to the terminal, print it pretty."`
}

var Repo = cmd.Sub{
	Name:  "repo",
	Alias: "r",
	Short: "Generate stats (json) for repositories.",
	Flags: &RepoFlags{},
	Args:  &RepoArgs{},
	Run:   RunRepo,
}

func init() {
	cmd.Register(&Repo)
}

func RunRepo(r *cmd.Root, c *cmd.Sub) {
	args := c.Args.(*RepoArgs)
	flags := c.Flags.(*RepoFlags)

	// a lookup of repo results by org
	results := map[string]github.RepoResult{}
	for _, repo := range args.Repos {
		results[repo] = github.GetRepoStats(repo)
	}

	// Parse into json
	var outJson []byte
	if flags.Pretty {
		outJson, _ = json.MarshalIndent(results, "", "    ")
	} else {
		outJson, _ = json.Marshal(results)
	}
	output := string(outJson)
	fmt.Printf(output)
}
