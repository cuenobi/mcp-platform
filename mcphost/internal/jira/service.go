package jira

import "os"

type Service struct {
	client Client
}

func NewService() *Service {
	// Use environment variable for Docker deployment, fallback to localhost
	addr := os.Getenv("MCP_SERVER_JIRA_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	return &Service{
		client: NewGRPCClient(addr),
	}
}

func (s *Service) Sync(project string) error {
	return s.client.Sync(project)
}

func (s *Service) CreateCard(project, prompt string) (string, error) {
	issueKey, err := s.client.CreateCard(project, prompt)
	if err != nil {
		return "", err
	}
	return issueKey, nil
}
