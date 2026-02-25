package generatemocks

import (
	"binh-swagger/cmd/swagger/commands/generate"
	"binh-swagger/cmd/swagger/commands/internal/spec"
)

type MockOrchestrator struct {
	FromAPIConfigFn func(*spec.APIConfig, generate.Config) error
}

func (m *MockOrchestrator) FromAPIConfig(cfg *spec.APIConfig, command generate.Config) error {
	if m.FromAPIConfigFn != nil {
		return m.FromAPIConfigFn(cfg, command)
	}

	return nil
}
