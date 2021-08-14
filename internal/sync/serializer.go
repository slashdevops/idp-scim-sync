package sync

type SyncSerializer interface {
	Decode(input []byte) (*GroupsResult, error)
	Encode(input *GroupsResult) ([]byte, error)
}
