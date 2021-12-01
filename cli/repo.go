package cli

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/vsoch/codestats/github"
	"io/ioutil"
)

// Args and flags
type RepoArgs struct {
	Repos []string `desc:"One or more GitHub repository names to parse."`
}
type RepoFlags struct {
	Pretty  bool   `long:"pretty" desc:"If printing to the terminal, print it pretty."`
	Outfile string `long:"outfile" desc:"Save output to file."`
	Config  string `long:"config" desc:"Provide a config to select metrics."`
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
	results := []github.RepoResult{}
	for _, repo := range args.Repos {
		results = append(results, github.GetRepoStats(repo, flags.Config))
	}

	// Parse into json
	var outJson []byte
	if flags.Pretty || flags.Outfile != "" {
		outJson, _ = json.MarshalIndent(results, "", "    ")
	} else {
		outJson, _ = json.Marshal(results)
	}
	if flags.Outfile != "" {
		_ = ioutil.WriteFile(flags.Outfile, outJson, 0644)
	} else {
		fmt.Printf(string(outJson))
	}
}
