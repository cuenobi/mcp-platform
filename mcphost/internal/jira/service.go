package jira

type Service struct {
    client Client
}

func NewService() *Service {
    return &Service{
        client: NewGRPCClient("localhost:50051"), // หรือดึงจาก env/config
    }
}

func (s *Service) Sync(project string) error {
    return s.client.Sync(project)
}