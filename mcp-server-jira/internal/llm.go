package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type IssueIdea struct {
	Title       string
	Description string
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func ReceivePrompt(prompt string) (string, error) {
	decisionPrompt := fmt.Sprintf(`You are a routing assistant. I will give you a message, and you need to decide how to handle it.

Answer "mcp" if the message requires any of these actions:
- Creating Jira cards/issues/tickets
- Updating Jira issues
- Syncing with Jira
- Any task management or project management actions
- Technical implementation tasks that need to be tracked

Answer "local" if the message is:
- Simple greetings (Hello, Hi, สวัสดี, etc.)
- General questions
- Casual conversation
- Requests for information only

User message: %s

Your answer (mcp or local):`, prompt)

	fmt.Printf("🔍 DEBUG: Sending decision prompt to Ollama: %s\n", decisionPrompt)

	response, err := callOllama(decisionPrompt)
	if err != nil {
		return "", err
	}

	fmt.Printf("🔍 DEBUG: Ollama response: '%s'\n", response)

	cleaned := strings.ToLower(strings.TrimSpace(response))
	fmt.Printf("🔍 DEBUG: Cleaned response: '%s'\n", cleaned)

	if strings.Contains(cleaned, "mcp") {
		fmt.Printf("🔍 DEBUG: Routing to MCP server\n")
		return talkToMCPServer(prompt)
	} else if strings.Contains(cleaned, "local") {
		fmt.Printf("🔍 DEBUG: Routing to local handler\n")
		return talkLocally(prompt)
	} else {
		fmt.Printf("🔍 DEBUG: No clear routing decision, defaulting to local\n")
		return talkLocally(prompt)
	}
}

func callOllama(prompt string) (string, error) {
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}

	jsonPayload, _ := json.Marshal(payload)

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func talkToMCPServer(prompt string) (string, error) {
	lowerPrompt := strings.ToLower(prompt)
	if strings.Contains(lowerPrompt, "create") && (strings.Contains(lowerPrompt, "card") || strings.Contains(lowerPrompt, "issue") || strings.Contains(lowerPrompt, "ticket") || strings.Contains(lowerPrompt, "jira")) {
		fmt.Println("Detected Jira card creation request via message command")

		projectKey := os.Getenv("JIRA_PROJECT_KEY")
		if projectKey == "" {
			projectKey = "PROJ"
		}

		issueIdea, err := GenerateIssueIdea(prompt)
		if err != nil {
			return "", fmt.Errorf("failed to generate issue idea: %w", err)
		}

		issueKey, err := CreateIssue(projectKey, issueIdea.Title, issueIdea.Description)
		if err != nil {
			return "", fmt.Errorf("failed to create Jira issue: %w", err)
		}

		return fmt.Sprintf("✅ Created Jira card: %s\nTitle: %s\nDescription: %s", issueKey, issueIdea.Title, issueIdea.Description), nil
	}

	return "Sending to MCP server: " + prompt, nil
}

func talkLocally(prompt string) (string, error) {
	return "Answer locally: " + prompt, nil
}

func GenerateIssueIdea(prompt string) (*IssueIdea, error) {
	fmt.Println("-------------- MCP Server Jira Generate Issue Idea ------------------")

	payload := map[string]interface{}{
		"model": "llama3",
		"prompt": prompt + "\n\nJIRA requirements:\n" +
			"1. Title must be less than 255 characters.\n" +
			"2. Description must be less than 1000 characters.\n" +
			"3. Summary must be less than 255 characters.\n\n" +
			"Please respond in this format:\nTitle: <your title here>\nDescription:\n<your description here>",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/generate", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to Ollama failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Bytes()
		var chunk struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}
		if err := json.Unmarshal(line, &chunk); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
		}
		fullResponse.WriteString(chunk.Response)
		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	content := fullResponse.String()

	title := extractTitle(content)
	description := extractDescription(content)

	title = sanitizeTitle(title)

	return &IssueIdea{
		Title:       title,
		Description: description,
	}, nil
}

func extractTitle(content string) string {
	re := regexp.MustCompile(`(?i)title:\s*(.+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	return "Untitled"
}

func extractDescription(content string) string {
	parts := strings.SplitN(content, "Description:", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return "No description provided"
}

func sanitizeTitle(title string) string {
	title = strings.ReplaceAll(title, "\n", " ")
	if len(title) > 255 {
		title = title[:255]
	}
	return title
}
