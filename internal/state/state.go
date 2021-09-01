package state

type State struct{}

func NewSyncState(name string) *State {
	return &State{}
}

func (s *State) GetName() string {
	return ""
}

func (s *State) isEmpty() bool {
	return true
}
