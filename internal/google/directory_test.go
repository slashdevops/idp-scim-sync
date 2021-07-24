package google

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	admin "google.golang.org/api/admin/directory/v1"
)

type GroupsServiceActions interface {
	List() *admin.GroupsListCall
}

type GroupsListCallActions interface {
	Customer(customer string) *admin.GroupsListCall
	Query(query string) *admin.GroupsListCall
	Pages(ctx context.Context, f func(*admin.Groups) error) error
}

type mockGroupsServiceActions struct {
	mock.Mock
	GroupsServiceActions
}

func (m *mockGroupsServiceActions) List() *admin.GroupsListCall {
	r := m.Called()
	return r.Get(0).(*admin.GroupsListCall)
}

type mockGroupsListCallActions struct {
	mock.Mock
	GroupsListCallActions
}

func (m *mockGroupsListCallActions) Customer(customer string) *admin.GroupsListCall {
	r := m.Called()
	return r.Get(0).(*admin.GroupsListCall)
}

func (m *mockGroupsListCallActions) Query(query string) *admin.GroupsListCall {
	r := m.Called()
	return r.Get(0).(*admin.GroupsListCall)
}

func (m *mockGroupsListCallActions) Pages(ctx context.Context, f func(*admin.Groups) error) error {
	r := m.Called()
	return r.Error(1)
}

func NewMockService() *admin.Service {

	ags := &admin.GroupsService{}

	svc := &admin.Service{
		Groups: ags,
	}

	return svc
}

func TestNewDirectoryService(t *testing.T) {

	t.Run("New Client with parameters", func(t *testing.T) {

		ctx := context.TODO()
		userEmail := "mock-email@mock-project.iam.gserviceaccount.com"
		serviceAccountFile := "../testdata/service_account.json"

		serviceAccount, err := ioutil.ReadFile(serviceAccountFile)
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}

		client, err := NewDirectoryService(ctx, userEmail, serviceAccount)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("New Client without right parameters", func(t *testing.T) {

		ctx := context.TODO()
		userEmail := ""
		serviceAccount := []byte("")

		client, err := NewDirectoryService(ctx, userEmail, serviceAccount)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

// func Test_directory_ListGroups(t *testing.T) {

// 	t.Run("ListGroups", func(t *testing.T) {
// 		ctx := context.TODO()
// 		svc := NewMockService()
// 		query := []string{"name:mock1"}

// 		gsa := &mockGroupsServiceActions{}
// 		gsa.On("List").Return(nil)

// 		glca := mockGroupsListCallActions{}
// 		glca.On("Customer", "my_customer").Return(&admin.GroupsListCall{})
// 		glca.On("Query", query).Return(&admin.GroupsListCall{})
// 		glca.On("Pages", ctx, mock.AnythingOfType("func(*admin.Groups) error")).Return(nil)

// 		d := &directory{
// 			ctx: ctx,
// 			svc: svc,
// 		}

// 		grps, err := d.ListGroups(query)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, grps)

// 		assert.Equal(t, "mock1", grps[0].Name)
// 	})
// }
