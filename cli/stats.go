package cli

import (
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/vsoch/org-stats/github"
)

// Args and flags
type StatsArgs struct {
	Orgs []string `desc:"One or more GitHub organization names to parse."`
}
type StatsFlags struct{}

var Stats = cmd.Sub{
	Name:  "stats",
	Alias: "s",
	Short: "Generate stats (json) for org projects.",
	Flags: &StatsFlags{},
	Args:  &StatsArgs{},
	Run:   RunStats,
}

func init() {
	cmd.Register(&Stats)
}

func RunStats(r *cmd.Root, c *cmd.Sub) {
	args := c.Args.(*StatsArgs)

	// a lookup of repo results by org
	results := map[string][]github.RepoResult{}
	for _, org := range args.Orgs {
		results[org] = github.GetOrgStats(org)
	}
	fmt.Println(results)
}
