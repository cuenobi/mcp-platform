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

var prompt string

var jiraCreateCmd = &cobra.Command{
	Use:   "create-card",
	Short: "Create Jira issue from prompt",
	Run: func(cmd *cobra.Command, args []string) {
		svc := jira.NewService()
		issueKey, err := svc.CreateCard(project, prompt)
		if err != nil {
			fmt.Printf("error creating card: %v\n", err)
			return
		}
		fmt.Printf("Created issue: %s\n", issueKey)
	},
}

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Send message to MCP server",
	Run: func(cmd *cobra.Command, args []string) {
		svc := jira.NewService()
		message, err := svc.Message(prompt)
		if err != nil {
			fmt.Printf("error sending message: %v\n", err)
			return
		}
		fmt.Printf("Message: %s\n", message)
	},
}

func init() {
	jiraSyncCmd.Flags().StringVarP(&project, "project", "p", "", "Jira project key")
	_ = jiraSyncCmd.MarkFlagRequired("project")
	jiraCmd.AddCommand(jiraSyncCmd)

	jiraCreateCmd.Flags().StringVarP(&project, "project", "p", "", "Jira project key")
	jiraCreateCmd.Flags().StringVarP(&prompt, "prompt", "", "", "Prompt to generate issue")
	_ = jiraCreateCmd.MarkFlagRequired("project")
	_ = jiraCreateCmd.MarkFlagRequired("prompt")

	messageCmd.Flags().StringVarP(&prompt, "prompt", "", "", "Prompt to send to MCP server")
	_ = messageCmd.MarkFlagRequired("prompt")

	jiraCmd.AddCommand(jiraCreateCmd)
	jiraCmd.AddCommand(messageCmd)
	rootCmd.AddCommand(jiraCmd)
}
