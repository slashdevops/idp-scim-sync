package state

type SyncState struct {
	name string
}

func NewSyncState(name string) *SyncState {
	return &SyncState{
		name: name,
	}
}

func (s *SyncState) GetName() string {
	return s.name
}
