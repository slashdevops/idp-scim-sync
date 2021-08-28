// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/slashdevops/idp-scim-sync/internal/model"
)

// MockSyncRepository is a mock of SyncRepository interface.
type MockSyncRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSyncRepositoryMockRecorder
}

// MockSyncRepositoryMockRecorder is the mock recorder for MockSyncRepository.
type MockSyncRepositoryMockRecorder struct {
	mock *MockSyncRepository
}

// NewMockSyncRepository creates a new mock instance.
func NewMockSyncRepository(ctrl *gomock.Controller) *MockSyncRepository {
	mock := &MockSyncRepository{ctrl: ctrl}
	mock.recorder = &MockSyncRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncRepository) EXPECT() *MockSyncRepositoryMockRecorder {
	return m.recorder
}

// GetGroups mocks base method.
func (m *MockSyncRepository) GetGroups(place string) (*model.GroupsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", place)
	ret0, _ := ret[0].(*model.GroupsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockSyncRepositoryMockRecorder) GetGroups(place interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockSyncRepository)(nil).GetGroups), place)
}

// GetGroupsMembers mocks base method.
func (m *MockSyncRepository) GetGroupsMembers(place string) (*model.GroupsMembersResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsMembers", place)
	ret0, _ := ret[0].(*model.GroupsMembersResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsMembers indicates an expected call of GetGroupsMembers.
func (mr *MockSyncRepositoryMockRecorder) GetGroupsMembers(place interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsMembers", reflect.TypeOf((*MockSyncRepository)(nil).GetGroupsMembers), place)
}

// GetState mocks base method.
func (m *MockSyncRepository) GetState() (model.SyncState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState")
	ret0, _ := ret[0].(model.SyncState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockSyncRepositoryMockRecorder) GetState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockSyncRepository)(nil).GetState))
}

// GetUsers mocks base method.
func (m *MockSyncRepository) GetUsers(place string) (*model.UsersResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", place)
	ret0, _ := ret[0].(*model.UsersResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockSyncRepositoryMockRecorder) GetUsers(place interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockSyncRepository)(nil).GetUsers), place)
}

// StoreGroups mocks base method.
func (m *MockSyncRepository) StoreGroups(gr *model.GroupsResult) (model.StoreGroupsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreGroups", gr)
	ret0, _ := ret[0].(model.StoreGroupsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreGroups indicates an expected call of StoreGroups.
func (mr *MockSyncRepositoryMockRecorder) StoreGroups(gr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreGroups", reflect.TypeOf((*MockSyncRepository)(nil).StoreGroups), gr)
}

// StoreGroupsMembers mocks base method.
func (m *MockSyncRepository) StoreGroupsMembers(gr *model.GroupsMembersResult) (model.StoreGroupsMembersResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreGroupsMembers", gr)
	ret0, _ := ret[0].(model.StoreGroupsMembersResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreGroupsMembers indicates an expected call of StoreGroupsMembers.
func (mr *MockSyncRepositoryMockRecorder) StoreGroupsMembers(gr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreGroupsMembers", reflect.TypeOf((*MockSyncRepository)(nil).StoreGroupsMembers), gr)
}

// StoreState mocks base method.
func (m *MockSyncRepository) StoreState(state *model.SyncState) (model.StoreStateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreState", state)
	ret0, _ := ret[0].(model.StoreStateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreState indicates an expected call of StoreState.
func (mr *MockSyncRepositoryMockRecorder) StoreState(state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreState", reflect.TypeOf((*MockSyncRepository)(nil).StoreState), state)
}

// StoreUsers mocks base method.
func (m *MockSyncRepository) StoreUsers(ur *model.UsersResult) (model.StoreUsersResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUsers", ur)
	ret0, _ := ret[0].(model.StoreUsersResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreUsers indicates an expected call of StoreUsers.
func (mr *MockSyncRepositoryMockRecorder) StoreUsers(ur interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUsers", reflect.TypeOf((*MockSyncRepository)(nil).StoreUsers), ur)
}