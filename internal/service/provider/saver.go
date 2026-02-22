package provider

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (s *Provider) Get() error {
	return nil
}
