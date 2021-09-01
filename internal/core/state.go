package core

type SyncState interface {
	GetName() string
	IsEmpty() bool
}
