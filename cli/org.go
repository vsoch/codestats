package cli

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/vsoch/codestats/github"
)

// Args and flags
type OrgArgs struct {
	Orgs []string `desc:"One or more GitHub organization names to parse."`
}
type OrgFlags struct {
	Pretty  bool   `long:"pretty" desc:"If printing to the terminal, print it pretty."`
	Pattern string `long:"pattern" desc:"Only include repos that match this regular expression."`
}

var Org = cmd.Sub{
	Name:  "org",
	Alias: "o",
	Short: "Generate stats (json) for org projects.",
	Flags: &OrgFlags{},
	Args:  &OrgArgs{},
	Run:   RunOrg,
}

func init() {
	cmd.Register(&Org)
}

func RunOrg(r *cmd.Root, c *cmd.Sub) {
	args := c.Args.(*OrgArgs)
	flags := c.Flags.(*OrgFlags)

	// a lookup of repo results by org
	results := map[string][]github.RepoResult{}
	for _, org := range args.Orgs {
		results[org] = github.GetOrgStats(org, flags.Pattern)
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
