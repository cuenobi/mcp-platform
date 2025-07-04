package cmd

import (
	"fmt"

	"github.com/cuenobi/mcp-platform/mcphost/internal/jira"
	"github.com/spf13/cobra"
)

var project string

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Interact with Jira MCP server",
}

var jiraSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Jira issues",
	Run: func(cmd *cobra.Command, args []string) {
		svc := jira.NewService()
		if err := svc.Sync(project); err != nil {
			fmt.Printf("error syncing jira: %v\n", err)
		}
	},
}

func init() {
	jiraSyncCmd.Flags().StringVarP(&project, "project", "p", "", "Jira project key")
	_ = jiraSyncCmd.MarkFlagRequired("project")
	jiraCmd.AddCommand(jiraSyncCmd)
	rootCmd.AddCommand(jiraCmd)
}
