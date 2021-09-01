package disk

import (
	"sync"
)

type DiskRepository struct {
	mu   *sync.RWMutex
	path string
}

func NewDiskRepository(path string) *DiskRepository {
	return &DiskRepository{}
}
