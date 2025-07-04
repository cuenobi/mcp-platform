package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/cuenobi/mcp-platform/mcphost/internal/jira"
)

var project string

var jiraCmd = &cobra.Command{
    Use:   "jira",
    Short: "Interact with Jira MCP server",
    Run: func(cmd *cobra.Command, args []string) {
        svc := jira.NewService()
        if err := svc.Sync(project); err != nil {
            fmt.Printf("error syncing jira: %v\n", err)
        }
    },
}

func init() {
    jiraCmd.Flags().StringVarP(&project, "project", "p", "", "Jira project key")
    _ = jiraCmd.MarkFlagRequired("project")
}