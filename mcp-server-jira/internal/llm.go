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

// func GenerateIssueIdea(prompt string) (*IssueIdea, error) {
// 	fmt.Println("Generating issue idea with prompt:", prompt)
// 	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

// 	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
// 		Model: openai.GPT3Dot5Turbo,
// 		Messages: []openai.ChatCompletionMessage{
// 			{Role: "system", Content: "You are a project manager who is good at writing Jira issues"},
// 			{Role: "user", Content: fmt.Sprintf("Help write Jira issue from command: %s", prompt)},
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	content := resp.Choices[0].Message.Content

// 	return &IssueIdea{
// 		Title:       content,
// 		Description: content,
// 	}, nil
// }

func GenerateIssueIdea(prompt string) (*IssueIdea, error) {
	fmt.Println("Generating issue idea with prompt (Ollama):", prompt)

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

	// Use environment variable for Docker deployment, fallback to localhost
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

	// แยก Title กับ Description ออกจาก content
	title := extractTitle(content)
	description := extractDescription(content)

	// ตัดความยาว title ไม่เกิน 255 ตัวอักษร และลบ newline
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
