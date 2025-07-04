package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateIssue(projectKey, title, description string) (string, error) {
	fmt.Println("-------------- mcp-server-jira create issue start ------------------")
	fmt.Println("Creating issue with project:", projectKey, "and title:", title, "and description:", description)
	jiraBaseURL := os.Getenv("JIRA_BASE_URL")
	jiraEmail := os.Getenv("JIRA_EMAIL")
	jiraToken := os.Getenv("JIRA_API_TOKEN")

	if jiraBaseURL == "" || jiraEmail == "" || jiraToken == "" {
		return "", fmt.Errorf("missing required Jira credentials or URL in environment variables")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", jiraEmail, jiraToken)))

	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project":     map[string]string{"key": projectKey},
			"summary":     title,
			"description": description,
			"issuetype":   map[string]string{"name": "Task"},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal issue payload: %w", err)
	}
	// Debug: log the JSON payload before sending the request
	fmt.Println("Jira CreateIssue JSON payload:", string(body))

	req, err := http.NewRequest("POST", jiraBaseURL+"/rest/api/2/issue", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Read and log full response body for debugging
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Jira API returned status: %s\nResponse body: %s\n", resp.Status, string(respBody))
		return "", fmt.Errorf("Jira API returned status: %s", resp.Status)
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("-------------- mcp-server-jira create issue end ------------------")
	return result.Key, nil
}
