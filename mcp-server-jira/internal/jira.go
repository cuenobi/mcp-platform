package internal

import "fmt"

func CreateIssue(projectKey, title, description string) (string, error) {
	fmt.Printf("Creating Jira issue in project %s with title: %s\n", projectKey, title)
	return "JIRA-1234", nil
}
