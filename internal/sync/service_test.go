package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ProviderService mock
type mockProviderService struct{}

func (m *mockProviderService) GetGoups(*context.Context, []string) (*GroupResult, error) {
	return &GroupResult{}, nil
}

func (m *mockProviderService) GetUsers(*context.Context, []string) (*UserResult, error) {
	return &UserResult{}, nil
}

func (m *mockProviderService) GetGroupsMembers(*context.Context, *GroupResult) (*MemberResult, error) {
	return &MemberResult{}, nil
}

func (m *mockProviderService) GetUsersFromGroupsMembers(*context.Context, []string, *MemberResult) (*UserResult, error) {
	return &UserResult{}, nil
}

// SCIMService mock
type mockSCIMService struct{}

func (m *mockSCIMService) GetGroups(*context.Context, []string) (*GroupResult, error) {
	return &GroupResult{}, nil
}

func (m *mockSCIMService) GetUsers(*context.Context, []string) (*UserResult, error) {
	return &UserResult{}, nil
}

func (m *mockSCIMService) CreateOrUpdateGroups(*context.Context, *GroupResult) error {
	return nil
}

func (m *mockSCIMService) CreateOrUpdateUsers(*context.Context, *UserResult) error {
	return nil
}

func (m *mockSCIMService) DeleteGroups(*context.Context, *GroupResult) error {
	return nil
}

func (m *mockSCIMService) DeleteUsers(*context.Context, *UserResult) error {
	return nil
}

func TestNewSyncService(t *testing.T) {

	t.Run("New Service with parameters", func(t *testing.T) {

		ctx := context.TODO()
		mockProviderService := &mockProviderService{}
		mockSCIMService := &mockSCIMService{}

		svc, err := NewSyncService(&ctx, mockProviderService, mockSCIMService)

		assert.NoError(t, err)
		assert.NotNil(t, svc)

	})

	t.Run("New Service without parameters", func(t *testing.T) {
		ctx := context.TODO()

		svc, err := NewSyncService(&ctx, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, svc)

	})
}
