package service

import "github.com/Tyz3/nymgraph/internal/state"

type ConfigService struct {
	state *state.State
}

func NewConfigService(state *state.State) *ConfigService {
	return &ConfigService{
		state: state,
	}
}

func (s *ConfigService) UseProxy(b bool) {
	s.state.GetConfig().UseProxy = b
	s.state.GetConfig().Save()
}

func (s *ConfigService) SetProxy(proxy string) {
	s.state.GetConfig().Proxy = proxy
	s.state.GetConfig().Save()
}

func (s *ConfigService) UsingProxy() bool {
	return s.state.GetConfig().UseProxy
}

func (s *ConfigService) GetProxy() string {
	return s.state.GetConfig().Proxy
}

func (s *ConfigService) DeleteHistoryAfterQuit() bool {
	return s.state.GetConfig().DeleteHistoryAfterQuit
}

func (s *ConfigService) SetDeleteHistoryAfterQuit(b bool) {
	s.state.GetConfig().DeleteHistoryAfterQuit = b
	s.state.GetConfig().Save()
}
