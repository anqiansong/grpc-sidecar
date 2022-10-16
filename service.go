package proxy

type Service struct {
	Name    string
	Address string
}

func (s *Service) Ping() error {
	panic("implement me")
}
