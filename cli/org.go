package cli

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/vsoch/codestats/github"
	"io/ioutil"
	"strings"
)

// Args and flags
type OrgArgs struct {
	Orgs []string `desc:"One or more GitHub organization names to parse."`
}
type OrgFlags struct {
	Pretty      bool   `long:"pretty" desc:"If printing to the terminal, print it pretty."`
	Pattern     string `long:"pattern" desc:"Only include repos that match this regular expression."`
	SkipPattern string `long:"skip" desc:"Skip repositories that match this pattern."`
	Outfile     string `long:"outfile" desc:"Save output to file."`
	Config      string `long:"config" desc:"Provide a config to select metrics."`
	Metric      string `long:"metric" desc:"A single metrics to provide on command line (overrides config)."`
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

	// Split metric by comma
	metrics := strings.Split(flags.Metric, ",")
	if metrics[0] == "" {
		metrics = []string{}
	}

	// a lookup of repo results by org
	results := []github.RepoResult{}
	for _, org := range args.Orgs {
		results = append(results, github.GetOrgStats(org, flags.Pattern, flags.SkipPattern, flags.Config, metrics)...)
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
