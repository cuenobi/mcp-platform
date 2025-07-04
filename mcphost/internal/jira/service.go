package jira

type Service struct {
	client Client
}

func NewService() *Service {
	return &Service{
		client: NewGRPCClient("localhost:50051"), // TODO: get from config
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
