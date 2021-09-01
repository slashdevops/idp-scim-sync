package state

type State struct {
	Name          string `json:"name"`
	SchemaVersion string `json:"version"`
	HashCode      string `json:"hashCode"`
	Resources     StateResources
}

type StateResources struct {
	StoreGroupsResult
	StoreUsersResult
	StoreGroupsUsersResult
}

type StoreGroupsResult struct {
	Location string
}

type StoreUsersResult struct {
	Location string
}

type StoreGroupsUsersResult struct {
	Location string
}

type StoreStateResult struct {
	Location string
}

func NewState(name string) *State {
	return &State{
		Name:          name,
		SchemaVersion: "1.0",
	}
}

func (s *State) GetName() string {
	return s.Name
}

func (s *State) Empty() bool {
	if s != nil {
		return false
	}

	if s.Resources.StoreGroupsResult.Location != "" {
		return false
	}
	if s.Resources.StoreGroupsUsersResult.Location != "" {
		return false
	}
	if s.Resources.StoreUsersResult.Location != "" {
		return false
	}

	return true
}

func (s *State) Build(groups *StoreGroupsResult, groupsUsers *StoreGroupsUsersResult, users *StoreUsersResult) (*State, error) {
	return &State{}, nil
}
