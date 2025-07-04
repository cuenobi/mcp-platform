package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		"model":  "llama3",
		"prompt": prompt + "\n\nJIRA requirements:\n1. Title must be less than 255 characters.\n2. Description must be less than 1000 characters.\n3. Summary must be less than 255 characters.\n\nPlease respond in this format:\nTitle: <your title here>\nDescription:\n<your description here>",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/generate", bytes.NewBuffer(payloadBytes))
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
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
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

	return &IssueIdea{
		Title:       content,
		Description: content,
	}, nil
}
